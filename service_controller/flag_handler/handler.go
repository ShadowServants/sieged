package flaghandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"hackforces/libs/flagdata"
	"hackforces/libs/flagresponse"
	"hackforces/libs/helpers"
	"hackforces/libs/round_storage"
	"hackforces/libs/statusstorage"
	"hackforces/libs/storage"
	"hackforces/service_controller/flag_handler/flagstorage"
	"math"
	"strconv"
	"sync"
)

//func handle(data string)

const SelfFLagMessage = "self"
const BadFlagMessage = "invalid"
const AlreadySubmitMessage = "already_submitted"
const TeamNotFoundMessage = "team_not_found"
const FlagTooOldMessage = "too_old"

type TeamData struct {
	id     int
	points Points
	mu     sync.Mutex
}

type Points struct {
	Plus   int     `json:"plus"`
	Minus  int     `json:"minus"`
	Points float64 `json:"points"`
}

type TeamRequest struct {
	Flag string `json:"flag"`
	Team int    `json:"team"`
}

func DumpsTeamRequest(tr *TeamRequest) (string, error) {
	byt, err := json.Marshal(tr)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

func LoadsTeamRequest(s string) (*TeamRequest, error) {
	tr := TeamRequest{}
	if err := json.Unmarshal([]byte(s), &tr); err != nil {
		return nil, errors.New("json_unmarshal_error")
	}
	return &tr, nil
}

func NewTeamData(id int, points Points) *TeamData {
	return &TeamData{id: id, points: points, mu: sync.Mutex{}}
}

type FlagHandler struct {
	Teams         map[int]*TeamData
	Flags         *flagstorage.FlagStorage
	Points        *PointsStorage
	TeamFlagsSet  storage.KeySet
	RoundSt       *round_storage.RoundStorage
	StatusStorage *statusstorage.StatusStorage
	RoundCached   bool
	CurrentRound  int
	RoundDelta    int
	TeamNum       int
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

	fh.RoundSt = &round_storage.RoundStorage{st}
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

func (fh *FlagHandler) CacheRound(callback func()) {
	fh.RoundCached = true
}

func (fh *FlagHandler) calcDelta(attackerPoints float64, victimPoints float64) float64 {
	if attackerPoints < victimPoints {
		return float64(fh.TeamNum)
	}
	ap := float64(attackerPoints)
	vp := float64(victimPoints)
	fmt.Println(ap, vp)
	logattacker := math.Max(math.Log(ap+1), 1)
	logvictim := math.Max(math.Log(vp+1), 1)
	delta := logvictim / logattacker
	deltaPointsF := math.Exp(math.Pow(math.Log(float64(fh.TeamNum)), delta))
	fmt.Println(deltaPointsF)
	//delta_points := int(delta_points_f)
	//delta_points := int(delta * float64(fh.TeamNum))
	return deltaPointsF

}

func (fh *FlagHandler) calc(att int, vict int) float64 {
	min := helpers.MinInt(att, vict)
	max := helpers.MaxInt(att, vict)
	attacker, _ := fh.GetTeamDataById(att)
	victim, _ := fh.GetTeamDataById(vict)

	fh.Teams[min].mu.Lock()
	fh.Teams[max].mu.Lock()
	defer fh.Teams[min].mu.Unlock()
	defer fh.Teams[max].mu.Unlock()
	delta := fh.calcDelta(attacker.points.Points, victim.points.Points)
	attacker.points.Points += delta
	attacker.points.Plus += 1
	victim.points.Minus += 1
	victim.points.Points -= math.Min(victim.points.Points, delta)
	//victim.points.Points = helpers.MaxInt(victim.points.Points-delta,0)
	go fh.StoreData(*attacker, *victim)
	return delta

}

func (fh *FlagHandler) CheckFlag(flag string) *flagdata.FlagData {
	fmt.Println("Try to send ", flag)
	if flag, err := fh.Flags.GetFlagData(flag); err == nil {
		return flag
	} else {
		fmt.Println(err.Error())
		return nil
	}

}

func (fh *FlagHandler) StoreData(teams ...TeamData) {
	for _, team := range teams {
		fh.Points.SetPoints(strconv.Itoa(team.id), &team.points)
	}
}

func (fh *FlagHandler) GetTeamDataById(id int) (*TeamData, error) {
	if data, ok := fh.Teams[id]; ok {
		return data, nil
	}
	if pts, err := fh.Points.GetPoints(strconv.Itoa(id)); err == nil {
		td := NewTeamData(id, *pts)
		fh.Teams[id] = td
		return td, nil
	}
	return nil, errors.New("team_not_found")
}

func (fh *FlagHandler) SetRoundCached(cached bool) *FlagHandler {
	fh.RoundCached = cached
	return fh
}

func (fh *FlagHandler) GetCurrentRound() int {
	if fh.RoundCached {
		return fh.CurrentRound
	} else {
		return fh.RoundSt.GetRound()
	}
}

func (fh *FlagHandler) ValidateFlag(tr *TeamRequest) (bool, string) {
	if _, err := fh.GetTeamDataById(tr.Team); err != nil {
		return false, TeamNotFoundMessage
	}
	exist := fh.TeamFlagsSet.Check(strconv.Itoa(tr.Team), tr.Flag)
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
		return false, FlagTooOldMessage
	}
	//if fh.StatusStorage.GetStatus(tr.Team,fh.GetCurrentRound()) != "Up" {
	//	return false,BadTeamStatusMessage
	//}
	//fmt.Println(data,err.Error())
	//return false, BadTeamStatusMessage
	//}
	return true, ""

}

func (fh *FlagHandler) SetCaptured(tr *TeamRequest) {
	fh.TeamFlagsSet.Add(strconv.Itoa(tr.Team), tr.Flag)
}

func (fh *FlagHandler) HandleRequest(s string) string {
	teamRequest, err := LoadsTeamRequest(s)
	if err != nil {
		return "Bad request"
	}
	response := flagresponse.NewStealResponse().SetInitiator(teamRequest.Team)
	if ok, responseText := fh.ValidateFlag(teamRequest); !ok {
		response.SetReason(responseText).SetSuccessful(false).SetTarget(-1)
		resp, _ := flagresponse.DumpHandlerResponse(response)
		return resp
	}
	victim := fh.CheckFlag(teamRequest.Flag).Team
	delta := fh.calc(teamRequest.Team, victim)
	fh.SetCaptured(teamRequest)
	response.SetSuccessful(true).SetTarget(victim).SetDelta(delta)

	resp, _ := flagresponse.DumpHandlerResponse(response)
	return resp

}
