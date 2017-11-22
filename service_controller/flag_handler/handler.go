package flaghandler

//import "github.com/johnnovikov/hackforces/back/service_controller/libs"
import (
	//"github.com/streadway/amqp"
	//"github.com/jnovikov/hackforces/back/libs/rpc"
	//redigo "github.com/garyburd/redigo/redis"
	//"github.com/streadway/amqp"
	"../../libs/helpers"
	//"github.com/jnovikov/hackforces/libs/helpers"
	"sync"
	//"os"
	//"os/signal"
	//"syscall"
	//"log"
	"../../libs/flagresponse"
	//"github.com/jnovikov/hackforces/libs/flagresponse"
	//"../../libs/flagresponse"
	//"github.com/jnovikov/hackforces/back/libs/flagresponse"
	"../../libs/storage"

	//"github.com/jnovikov/hackforces/libs/storage"
	"strconv"
	"math"
	"encoding/json"
	"errors"
	//"github.com/jnovikov/hackforces/libs/flagstorage"
	//"github.com/jnovikov/hackforces/libs/flagdata"
	//"github.com/jnovikov/hackforces/libs/statusstorage"
	"../../libs/statusstorage"
	"../../libs/flagdata"
	"./flagstorage"
)


//func handle(data string)

const SelfFLagMessage = "self"
const BadFlagMessage = "invalid"
const AlreadySubmitMessage = "already_submitted"
const TeamNotFoundMessage = "team_not_found"
const FlagTooOldMessage = "too_old"
const BadTeamStatusMessage  = "not_ok"


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

func DumpsTeamRequest(tr *TeamRequest) (string,error) {
	byt, err := json.Marshal(tr)
	if err != nil {
		return "",err
	}
	return string(byt),nil
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
	StatusStorage *statusstorage.StatusStorage
	RoundCached bool
	CurrentRound int
	RoundDelta int
	//pool *redigo.Pool
}

func NewFlagHandler() *FlagHandler {
	f := new(FlagHandler)
	f.Teams = make(map[int]*TeamData)

	f.RoundCached = false
	f.RoundDelta = 5
	return f
}

func (fh *FlagHandler) SetRoundStorage(st storage.Storage) *FlagHandler {

	fh.RoundSt = &RoundStorage{st}
	return fh
}

func (fh *FlagHandler) SetPointStorage(st storage.Storage) *FlagHandler {
	fh.Points = &PointsStorage{st}
	return fh
}

func (fh *FlagHandler) SetFlagStorage(st storage.Storage) *FlagHandler {
	fh.Flags = flagstorage.NewFlagStorage(st)
	return fh
}

func (fh *FlagHandler) SetStatusStorage(st storage.Storage) *FlagHandler {
	fh.StatusStorage = statusstorage.NewStatusStorage(st)
	return fh
}

func (fh *FlagHandler) SetTeamFlagsSet(ks storage.KeySet) *FlagHandler {
	fh.TeamFlagsSet = ks
	return fh
}

func (fh *FlagHandler) CacheRound(callback func ()) {
	fh.RoundCached = true
}

func (fh *FlagHandler) calcDelta(attacker_points int,victim_points int) int {
	ap := math.Max(1.0,float64(attacker_points + 1))
	vp := math.Max(1.0,float64(victim_points + 1))
	if ap > vp {
		return int(math.Exp(math.Log(15.0) * (15.0 - vp) / (15.0 - ap)))
	}
	logattacker := math.Log2(ap) + 1
        logvictim := math.Log2(vp) + 1
        delta := logvictim / logattacker
        delta_points := int(delta * 15)
	return delta_points

}


func (fh *FlagHandler) calc(att int, vict int) int{
	min := helpers.MinInt(att,vict)
	max := helpers.MaxInt(att,vict)
	attacker,_ := fh.GetTeamDataById(att)
	victim,_ := fh.GetTeamDataById(vict)

	fh.Teams[min].mu.Lock()
	fh.Teams[max].mu.Lock()
	defer fh.Teams[min].mu.Unlock()
	defer fh.Teams[max].mu.Unlock()
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

func (fh *FlagHandler)  GetTeamDataById(id int) (*TeamData,error){
	if data,ok := fh.Teams[id]; ok {
		return data,nil
	}
	if pts, err := fh.Points.GetPoints(strconv.Itoa(id)); err == nil {
		td := NewTeamData(id,*pts)
		fh.Teams[id] = td
		return td,nil
	}
	return nil,errors.New("Cant get team")
}

func (fr *FlagHandler) SetRoundCached(cached bool) *FlagHandler {
	fr.RoundCached = cached
	return fr
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
	if fh.StatusStorage.GetStatus(tr.Team,fh.GetCurrentRound()) != "Up" {
		return false,BadTeamStatusMessage
	}
		//fmt.Println(data,err.Error())
	 	//return false, BadTeamStatusMessage
	//}
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
	response := flagresponse.NewStealResponse().SetInitiator(team_request.Team)
	if ok,response_text := fh.ValidateFlag(team_request); !ok {
		response.SetReason(response_text).SetSuccessful(false).SetTarget(-1)
		resp, _ := flagresponse.DumpHandlerResponse(response)
		return resp
	}
	victim := fh.CheckFlag(team_request.Flag).Team
	delta := fh.calc(team_request.Team,victim)
	fh.SetCaptured(team_request)
	response.SetSuccessful(true).SetTarget(victim).SetDelta(delta)

	resp, _ := flagresponse.DumpHandlerResponse(response)
	return resp

}
