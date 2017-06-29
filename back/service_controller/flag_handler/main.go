package main

//import "github.com/johnnovikov/hackforces/back/service_controller/libs"
import (
	//"github.com/streadway/amqp"
	"github.com/jnovikov/hackforces/back/libs/rpc"
	//redigo "github.com/garyburd/redigo/redis"
	"github.com/streadway/amqp"
	"github.com/jnovikov/hackforces/back/libs/helpers"
	"sync"
	"os"
	"os/signal"
	"syscall"
	"log"
	"github.com/jnovikov/hackforces/back/libs/storage"
	"strconv"
	"math"
	"encoding/json"
	"errors"
	"github.com/jnovikov/hackforces/back/libs/flagstorage"
	"fmt"
	"github.com/jnovikov/hackforces/back/libs/flagdata"
)


//func handle(data string)

const SelfFLagMessage = "Thats your own flag"
const BadFlagMessage = "No such flag"
const AlreadySubmitMessage = "You already submit this flag"
const TeamNotFoundMessage = "Team not found"
const FlagTooOldMessage = "Flag is too old"


type TeamData struct {
	id int
	points Points
	mu sync.Mutex
}

type Points struct {
	Plus int `json:"plus"`
	Minus int `json:"minus"`
	Points int `json:"points"`
}




type TeamRequest struct {
	Flag string  `json:"flag"`
	Team int  `json:"team"`

}


func LoadsTeamRequest(s string) (*TeamRequest, error){
	tr := TeamRequest{}
	if err := json.Unmarshal([]byte(s), &tr); err != nil {
		return nil,errors.New("Cant unmarshall json")
    }
	return &tr,nil
}

func NewTeamData(id int,points Points) *TeamData{
	return &TeamData{id:id,points:points,mu:sync.Mutex{}}
}



type FlagHandler struct {
	Teams map[int] *TeamData
	Flags *flagstorage.FlagStorage
	Points *PointsStorage
	TeamFlagsSet storage.KeySet
	RoundSt *RoundStorage
	NumOfTeams int
	RoundCached bool
	CurrentRound int
	RoundDelta int
	//pool *redigo.Pool
}

func (fh *FlagHandler) calcDelta(attacker_points int,victim_points int) int{
		ap := math.Max(1.0,float64(attacker_points + 1))
        vp := math.Max(1.0,float64(victim_points + 1))
        logattacker := math.Log2(ap) + 1
        logvictim := math.Log2(vp) + 1
        delta := logvictim / logattacker
        delta_points := int(delta * 15)
		return delta_points
}


func (fh *FlagHandler) calc(att int, vict int) int{
	min := helpers.MinInt(att,vict)
	max := helpers.MaxInt(att,vict)
	fh.Teams[min].mu.Lock()
	fh.Teams[max].mu.Lock()
	defer fh.Teams[min].mu.Unlock()
	defer fh.Teams[max].mu.Unlock()
	attacker := fh.Teams[att]
	victim := fh.Teams[vict]
	delta := fh.calcDelta(attacker.points.Points,victim.points.Points)
	attacker.points.Points += delta
	attacker.points.Plus += 1
	victim.points.Minus += 1
	victim.points.Points = helpers.MaxInt(victim.points.Points-delta,0)
	go fh.StoreData(*attacker,*victim)
	return delta

}

func (fh *FlagHandler) CheckFlag(flag string) *flagdata.FlagData {
	if flag,err := fh.Flags.GetFlagData(flag); err == nil {
		return flag
	} else {
		return nil
	}

}

func (fh *FlagHandler) StoreData(teams ...TeamData) {
	for _, team := range teams {
		fh.Points.SetPoints(strconv.Itoa(team.id),&team.points)
	}
}

func (fh *FlagHandler) Build(num int,basepoints int) {
	fh.NumOfTeams = num
	fh.Teams = make(map[int]*TeamData)
	for i:=1; i<=num;i++{
		tp, err := fh.Points.GetPoints(strconv.Itoa(i))
		if err != nil {
			tp = &Points{0,0,basepoints}
		}
		td := NewTeamData(i,*tp)
		fh.Teams[i] =  td//points
		fh.StoreData(*td)
	}
}

func (fh *FlagHandler) GetTeamDataById(id int) (*TeamData,error){
	if data,ok := fh.Teams[id]; ok {
		return data,nil
	}
	return nil,errors.New("Cant get team")
}

func (fh *FlagHandler) GetCurrentRound() int {
	if fh.RoundCached {
		return fh.CurrentRound
	} else {
		return fh.RoundSt.GetRound()
	}
}

func (fh *FlagHandler) ValidateFlag(tr *TeamRequest) (bool,string) {
	if _,err := fh.GetTeamDataById(tr.Team); err != nil {
		return false, TeamNotFoundMessage
	}
	exist := fh.TeamFlagsSet.Check(strconv.Itoa(tr.Team),tr.Flag)
	if exist {
		return false, AlreadySubmitMessage
	}
	flag := fh.CheckFlag(tr.Flag)
	if flag == nil {
		return false, BadFlagMessage
	}
	victim := flag.Team
	if victim == tr.Team {
		return false, SelfFLagMessage
	}
	if helpers.Abs(fh.GetCurrentRound()-flag.Round) >= fh.RoundDelta {
		return false,FlagTooOldMessage
	}

	return true,""


}

func (fh *FlagHandler) SetCaptured(tr *TeamRequest) {
	fh.TeamFlagsSet.Add(strconv.Itoa(tr.Team),tr.Flag)
}

func (fh *FlagHandler) HandleRequest(s string) string{
	team_request,err := LoadsTeamRequest(s)
	if err != nil {
		return "Bad request"
	}
	if ok,response := fh.ValidateFlag(team_request); !ok {
		return response
	}
	victim := fh.CheckFlag(team_request.Flag).Team
	delta := fh.calc(team_request.Team,victim)
	fh.SetCaptured(team_request)
	return fmt.Sprintf("Congrats. You captured %d points",delta)

}


func main()  {
	conn , err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.FailOnError(err,"Cant connect to rabbit")
	defer conn.Close()
	a_handler:= new(rpc.AckHandler)
	a_handler.Init()

	mq := new(rpc.RabbitMqRpc)
	defer mq.Close()
	mq.Connection = conn
	mq.Build("flags_rpc",1)
	mq.Handler = a_handler
	go mq.Handle()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Println(<-ch)

	mq.Close()
}