package round_handler

import (
	"strconv"
	"sync"
	"os/exec"
	"errors"
	"encoding/json"
	"fmt"
	"hackforces/libs/storage"
	"hackforces/libs/round_storage"
	"hackforces/service_controller/flag_handler"
	"hackforces/libs/statusstorage"
	"hackforces/libs/team_list"
)

type RoundHandler struct {
	Wg            sync.WaitGroup
	IpStorage     storage.Storage
	TeamStorage   storage.Storage
	Points        *flaghandler.PointsStorage
	St            *statusstorage.StatusStorage
	TeamIds       []int
	Rounds        *round_storage.RoundStorage
	CheckerName   string
	DefaultPoints int
}

type RoundRequest struct {
	round int
}

type RoundResponse struct {
	Responses []TeamResponse `json:"responses"`
}

type TeamResponse struct {
	TeamId        int                `json:"team_id"`
	StatusMessage string             `json:"status_message"`
	Status        string             `json:"status"`
	Points        flaghandler.Points `json:"points"`
	//flaghandler.Points
}

func LoadTeamResponse(s string) (*TeamResponse, error) {
	tr := new(TeamResponse)
	if err := json.Unmarshal([]byte(s), tr); err != nil {
		return nil, errors.New("Cant unmarshall json points")
	}
	return tr, nil
}

func (rh *RoundHandler) GetIpByTeam(teamId int) string {
	ip, err := rh.IpStorage.Get(strconv.Itoa(teamId))
	if err != nil {
		return ""
	}
	return ip
}

func (rh *RoundHandler) TestTeam(teamId int, round int, ch chan TeamResponse) {
	defer rh.Wg.Done()
	ip := rh.GetIpByTeam(teamId)

	if ip == "" {
		ch <- TeamResponse{teamId, "Team not found", "Down", flaghandler.Points{}}
		rh.St.SetStatus(teamId, round, "Down")
		return
	}

	cmd := exec.Command(rh.CheckerName, "-t", strconv.Itoa(teamId), "-r", strconv.Itoa(round), "--ip", ip, "-tl", "7")
	stdout, err := cmd.Output()
	pts, _ := rh.Points.GetPoints(strconv.Itoa(teamId))
	if err != nil {
		fmt.Println("Checker failed with errors", err.Error(), stdout)
		rh.St.SetStatus(teamId, round, "Down")
		ch <- TeamResponse{teamId, "Checker failed", "Down", *pts}
		return
	}
	tr, err := LoadTeamResponse(string(stdout))
	tr.Points = *pts
	if err != nil {
		rh.St.SetStatus(teamId, round, "Down")
		ch <- TeamResponse{teamId, "Cant load response", "Down", *pts}
		return
	}
	rh.St.SetStatus(teamId, round, tr.Status)
	ch <- *tr

}

func (rh *RoundHandler) GetTeams() []int {
	if len(rh.TeamIds) != 0 {
		return rh.TeamIds
	}
	data, err := rh.TeamStorage.Get("teams_id")
	if err != nil {
		return nil
	}
	teams, err := team_list.LoadsTeamList(data)
	if err != nil {
		return nil
	}
	fmt.Println(teams)
	rh.TeamIds = teams.Teams
	return teams.Teams
}

func (rh *RoundHandler) CheckTeams(roundNum int) string {
	teams := rh.GetTeams()
	teamsNum := len(teams)
	rh.Wg.Add(teamsNum)
	ch := make(chan TeamResponse, teamsNum)
	for _, v := range teams {
		go rh.TestTeam(v, roundNum, ch)
	}
	rh.Wg.Wait()
	close(ch)

	responses := make([]TeamResponse, 0)
	for v := range ch {
		responses = append(responses, v)
	}
	fmt.Println(responses)
	rr := RoundResponse{responses}
	byt, err := json.Marshal(&rr)
	if err != nil {
		return "Bad json"
	}
	return string(byt)
}

func (rh *RoundHandler) HandleRequest(a string) string {
	roundInt, err := strconv.Atoi(a)
	if err != nil {
		return "Bad int"
	}
	rh.Rounds.SetRound(roundInt)
	return rh.CheckTeams(roundInt)

}

func (rh *RoundHandler) LoadsTeamsIp(teams []team_list.TeamIP) {
	teamsId := make([]int, 0)
	for _, team := range teams {
		teamIdStr := strconv.Itoa(team.Id)
		rh.IpStorage.Set(teamIdStr, team.Server)
		pts, _ := rh.Points.GetPoints(teamIdStr)
		if pts == nil {
			rh.Points.SetPoints(teamIdStr, &flaghandler.Points{0, 0, float64(rh.DefaultPoints)})
		}
		teamsId = append(teamsId, team.Id)
	}
	rh.TeamIds = teamsId
	tl := &team_list.TeamList{Teams: teamsId}

	tlString, _ := team_list.DumpTeamList(tl)
	rh.TeamStorage.Set("teams_id", tlString)

}

func NewRoundHandler() *RoundHandler {
	rh := new(RoundHandler)
	rh.Wg = sync.WaitGroup{}
	rh.TeamIds = make([]int, 0)
	rh.DefaultPoints = 1700
	return rh
}
