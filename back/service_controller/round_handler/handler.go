package round_handler

import (
	"github.com/jnovikov/hackforces/back/libs/storage"
	"github.com/jnovikov/hackforces/back/libs/team_list"

	"strconv"
	"github.com/jnovikov/hackforces/back/service_controller/flag_handler"
	"github.com/jnovikov/hackforces/back/libs/statusstorage"
	"sync"
	"os/exec"
	"errors"
	"encoding/json"
	"fmt"
)

type RoundHandler struct {
	Wg sync.WaitGroup
	IpStorage storage.Storage
	TeamStorage storage.Storage
	St *statusstorage.StatusStorage
	TeamIds []int
	Rounds *flaghandler.RoundStorage
	CheckerName string

}


type RoundRequest struct{
	round int
}

type RoundResponse struct {
	Responses []TeamResponse `json:"responses"`
}

type TeamResponse struct {
	Team_id int `json:"team_id"`
	Status_message string `json:"status_message"`
	Status string `json:"status"`
}


func LoadTeamResponse(s string) (*TeamResponse,error) {
	tr := new(TeamResponse)
	if err := json.Unmarshal([]byte(s), tr); err != nil {
		return nil,errors.New("Cant unmarshall json points")
    }
	return tr,nil
}

func (rh *RoundHandler) GetIpByTeam(team_id int) string {
	ip,err := rh.IpStorage.Get(strconv.Itoa(team_id))
	if err != nil {
		return ""
	}
	return ip
}



func (rh *RoundHandler) TestTeam(team_id int,round int,ch chan TeamResponse) {
	defer rh.Wg.Done()
	ip := rh.GetIpByTeam(team_id)
	fmt.Println(ip)
	if ip == "" {
		ch <- TeamResponse{team_id,"Team not found","Down"}
		rh.St.SetStatus(team_id,round,"Down")
		return
	}
	cmd := exec.Command(rh.CheckerName,"-t",strconv.Itoa(team_id),"-r",strconv.Itoa(round),"--ip",ip,"-tl","3")
	stdout, err := cmd.Output()
	if err != nil {
		rh.St.SetStatus(team_id,round,"Down")
		ch <- TeamResponse{team_id,"Checker failed","Down"}
		return
	}
	tr, err := LoadTeamResponse(string(stdout))
	if err != nil  {
		rh.St.SetStatus(team_id,round,"Down")
		ch <- TeamResponse{team_id,"Cant load response","Down"}
		return
	}
	rh.St.SetStatus(team_id,round,tr.Status)
	ch <- *tr




}


func (rh *RoundHandler) GetTeams() []int{
	if len(rh.TeamIds) != 0 {
		return rh.TeamIds
	}
	data,err := rh.TeamStorage.Get("teams_id")
	if err != nil {
		return nil
	}
	teams,err := team_list.LoadsTeamList(data)
	if err != nil {
		return nil
	}
	fmt.Println(teams)
	rh.TeamIds = teams.Teams
	return teams.Teams
}


func (rh *RoundHandler) CheckTeams(round_num int) string {
	teams := rh.GetTeams()
	teams_num := len(teams)
	rh.Wg.Add(teams_num)
	ch := make(chan TeamResponse, teams_num)
	for _, v := range teams {
		go rh.TestTeam(v,round_num,ch)
	}
	rh.Wg.Wait()
	close(ch)

	responses := make([]TeamResponse,0)
	for v := range ch {
		responses = append(responses,v)
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
	round_int , err := strconv.Atoi(a)
	if err != nil {
		return "Bad int"
	}
	rh.Rounds.SetRound(round_int)
	return rh.CheckTeams(round_int)


}