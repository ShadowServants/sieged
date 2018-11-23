package rounds

import (
	"context"
	"encoding/json"
	"fmt"
	"sieged/internal/team"
	"sieged/internal/team/score"
	"sieged/internal/team/status"
	"sieged/pkg/storage"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const RequestTimeout = "7"
const HardTimeOut = 30 * time.Second

type Handler struct {
	Wg            sync.WaitGroup
	IpStorage     storage.Storage
	TeamStorage   storage.Storage
	ScoreStorage  *score.Storage
	St            *status.Storage
	TeamIds       []int
	Rounds        *Storage
	CheckerName   string
	DefaultPoints int
}

func (rh *Handler) GetIpByTeam(teamId int) string {
	ip, err := rh.IpStorage.Get(strconv.Itoa(teamId))
	if err != nil {
		return ""
	}
	return ip
}

func (rh *Handler) TestTeam(teamId int, round int, ch chan team.Status) {
	defer rh.Wg.Done()
	ip := rh.GetIpByTeam(teamId)

	if ip == "" {
		ch <- team.Status{
			TeamId:        teamId,
			StatusMessage: "Team not found",
			Status:        "Down",
		}
		rh.St.SetStatus(teamId, round, "Down")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), HardTimeOut)
	defer cancel()
	cmd := exec.CommandContext(ctx, rh.CheckerName, "-t", strconv.Itoa(teamId), "-r", strconv.Itoa(round), "--ip", ip, "-tl", RequestTimeout)
	stdout, err := cmd.Output()
	pts, _ := rh.ScoreStorage.GetPoints(strconv.Itoa(teamId))
	if err != nil {
		fmt.Println("Checker failed with errors", err.Error(), string(stdout))
		rh.St.SetStatus(teamId, round, "Down")
		ch <- team.Status{TeamId: teamId, StatusMessage: "Timeout", Status: "Down", Points: *pts}
		return
	}
	tr, err := status.Loads(string(stdout))
	tr.Points = *pts
	if err != nil {
		rh.St.SetStatus(teamId, round, "Down")
		ch <- team.Status{TeamId: teamId, StatusMessage: "Cant load response", Status: "Down", Points: *pts}
		return
	}
	rh.St.SetStatus(teamId, round, tr.Status)
	ch <- *tr

}

func (rh *Handler) GetTeams() []int {
	if len(rh.TeamIds) != 0 {
		return rh.TeamIds
	}
	data, err := rh.TeamStorage.Get("teams_id")
	if err != nil {
		return nil
	}
	teams, err := team.LoadList(data)
	if err != nil {
		return nil
	}
	fmt.Println(teams)
	rh.TeamIds = teams.Teams
	return teams.Teams
}

func (rh *Handler) CheckTeams(roundNum int) string {
	teams := rh.GetTeams()
	teamsNum := len(teams)
	rh.Wg.Add(teamsNum)

	ch := make(chan team.Status, teamsNum)
	for _, v := range teams {
		go rh.TestTeam(v, roundNum, ch)
	}
	rh.Wg.Wait()
	close(ch)

	responses := make([]team.Status, 0)
	for v := range ch {
		responses = append(responses, v)
	}
	fmt.Println(responses)
	rr := Response{Responses: responses}
	byt, err := json.Marshal(&rr)
	if err != nil {
		return "Bad json"
	}
	return string(byt)
}

func (rh *Handler) HandleRequest(a string) string {
	roundInt, err := strconv.Atoi(a)
	if err != nil {
		return "Bad int"
	}
	rh.Rounds.SetRound(roundInt)
	return rh.CheckTeams(roundInt)

}

func (rh *Handler) LoadsTeamsIp(teams []team.Address) {
	teamsId := make([]int, 0)
	for _, currentTeam := range teams {
		teamIdStr := strconv.Itoa(currentTeam.Id)
		rh.IpStorage.Set(teamIdStr, currentTeam.Server)
		pts, _ := rh.ScoreStorage.GetPoints(teamIdStr)
		if pts == nil {
			rh.ScoreStorage.SetPoints(teamIdStr, &team.Score{Points: float64(rh.DefaultPoints)})
		}
		teamsId = append(teamsId, currentTeam.Id)
	}
	rh.TeamIds = teamsId
	tl := &team.List{Teams: teamsId}

	tlString, _ := team.DumpList(tl)
	rh.TeamStorage.Set("teams_id", tlString)

}

func NewHandler() *Handler {
	rh := new(Handler)
	rh.Wg = sync.WaitGroup{}
	rh.TeamIds = make([]int, 0)
	rh.DefaultPoints = 1700
	return rh
}
