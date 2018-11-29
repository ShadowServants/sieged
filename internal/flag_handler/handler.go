package flaghandler

import (
	"errors"
	"log"
	"math"
	"sieged/internal/flags"
	"sieged/internal/rounds"
	"sieged/internal/team"
	"sieged/internal/team/score"
	"sieged/internal/team/status"
	"sieged/pkg/helpers"
	"sieged/pkg/storage"
	"strconv"
)

const SelfFLagMessage = "self"
const BadFlagMessage = "invalid"
const AlreadySubmitMessage = "already_submitted"
const TeamNotFoundMessage = "team_not_found"
const FlagTooOldMessage = "too_old"

type FlagHandler struct {
	Teams         map[int]*team.Data
	Flags         *flags.Storage
	Points        *score.Storage
	TeamFlagsSet  storage.KeySet
	RoundSt       *rounds.Storage
	StatusStorage *status.Storage
	RoundCached   bool
	CurrentRound  int
	RoundDelta    int
	TeamNum       int
	//pool *redigo.Pool
}

func NewFlagHandler() *FlagHandler {
	f := new(FlagHandler)
	f.Teams = make(map[int]*team.Data)

	f.RoundCached = false
	f.RoundDelta = 3
	return f
}

func (fh *FlagHandler) SetRoundStorage(st storage.Storage) *FlagHandler {

	fh.RoundSt = &rounds.Storage{St: st}
	return fh
}

func (fh *FlagHandler) SetPointStorage(st storage.Storage) *FlagHandler {
	fh.Points = &score.Storage{St: st}
	return fh
}

func (fh *FlagHandler) SetFlagStorage(st storage.Storage) *FlagHandler {
	fh.Flags = flags.NewStorage(st)
	return fh
}

func (fh *FlagHandler) SetStatusStorage(st storage.Storage) *FlagHandler {
	fh.StatusStorage = status.NewStorage(st)
	return fh
}

func (fh *FlagHandler) SetTeamFlagsSet(ks storage.KeySet) *FlagHandler {
	fh.TeamFlagsSet = ks
	return fh
}

func (fh *FlagHandler) CacheRound(callback func()) {
	fh.RoundCached = true
}

func (fh *FlagHandler) R(B, A, t float64) float64 {
	if A+B == 0 {
		return 0.0
	}

	f := func(x float64) float64 {
		return math.Sqrt(x) / 14
	}
	rf := func(y float64) float64 {
		return math.Pow(y*14, 2)
	}

	return (A+B)*f(rf(B/(A+B))+t) - B

}

func (fh *FlagHandler) calcDelta(attacker float64, defender float64) float64 {
	//if attackerPoints < victimPoints {
	//	return float64(fh.TeamNum)
	//}
	//ap := float64(attackerPoints)
	//vp := float64(victimPoints)
	//fmt.Println(ap, vp)
	//logAttacker := math.Max(math.Log(ap+1), 1)
	//logVictim := math.Max(math.Log(vp+1), 1)
	//delta := logVictim / logAttacker
	//deltaPointsF := math.Exp(math.Pow(math.Log(float64(fh.TeamNum)), delta))
	//fmt.Println(deltaPointsF)
	//delta_points := int(delta_points_f)
	//delta_points := int(delta * float64(fh.TeamNum))
	r := fh.R(attacker, defender, 1.0)
	if attacker > defender && attacker > 0 {
		r = defender * r / attacker
	}
	if attacker < 240 && defender < 240 {
		r = math.Max(r, 3*(math.Floor(defender/10)+1))
	}
	return r

}

func (fh *FlagHandler) calc(att int, vict int) float64 {
	min := helpers.MinInt(att, vict)
	max := helpers.MaxInt(att, vict)
	attacker, _ := fh.GetTeamDataById(att)
	victim, _ := fh.GetTeamDataById(vict)

	fh.Teams[min].Lock()
	fh.Teams[max].Lock()
	defer fh.Teams[min].Unlock()
	defer fh.Teams[max].Unlock()
	delta := fh.calcDelta(attacker.Score.Points, victim.Score.Points)
	attacker.Score.Plus += 1
	victim.Score.Minus += 1

	attacker.Score.Points += delta + delta*0.10
	victim.Score.Points = math.Max(victim.Score.Points-delta, 0)
	//victim.points.Score = helpers.MaxInt(victim.points.Score-delta,0)
	go fh.StoreData(*attacker, *victim)
	return delta

}

func (fh *FlagHandler) CheckFlag(flag string) *flags.Data {
	if flag, err := fh.Flags.GetData(flag); err == nil {
		return flag
	} else {
		log.Println(err.Error())
		return nil
	}

}

func (fh *FlagHandler) StoreData(teams ...team.Data) {
	for _, t := range teams {
		fh.Points.SetPoints(strconv.Itoa(t.Id), &t.Score)
	}
}

func (fh *FlagHandler) GetTeamDataById(id int) (*team.Data, error) {
	if data, ok := fh.Teams[id]; ok {
		return data, nil
	}
	if pts, err := fh.Points.GetPoints(strconv.Itoa(id)); err == nil {
		td := &team.Data{Id: id, Score: *pts}
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

func (fh *FlagHandler) ValidateFlag(tr *flags.Request) (bool, string) {
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
	return true, ""

}

func (fh *FlagHandler) SetCaptured(tr *flags.Request) {
	fh.TeamFlagsSet.Add(strconv.Itoa(tr.Team), tr.Flag)
}

func (fh *FlagHandler) HandleRequest(s string) string {
	teamRequest, err := flags.LoadsRequest(s)
	if err != nil {
		return "Bad request"
	}
	response := flags.StealResponse().SetInitiator(teamRequest.Team)
	if ok, responseText := fh.ValidateFlag(teamRequest); !ok {
		response.SetReason(responseText).SetSuccessful(false).SetTarget(-1)
		resp, _ := flags.DumpResponse(response)
		return resp
	}
	victim := fh.CheckFlag(teamRequest.Flag).Team
	delta := fh.calc(teamRequest.Team, victim)
	fh.SetCaptured(teamRequest)
	response.SetSuccessful(true).SetTarget(victim).SetDelta(delta)

	resp, _ := flags.DumpResponse(response)
	return resp

}
