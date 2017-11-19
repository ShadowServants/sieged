package init_handler

import (
	"github.com/jnovikov/hackforces/back/service_controller/flag_handler"
	"github.com/jnovikov/hackforces/back/libs/storage"
	"github.com/jnovikov/hackforces/back/libs/team_list"
	"strconv"
	"fmt"
)

type InitHandler struct {
	Ps *flaghandler.PointsStorage
	TeamStorage storage.Storage

}



func (ih *InitHandler) HandleRequest(data string) string {
	ih.TeamStorage.Set("teams_id",data)
	teams,err := team_list.LoadsTeamList(data)
	if err != nil {
		return "BAD"
	}
	ih.TeamStorage.Set("team_num",strconv.Itoa(len(teams.Teams)))
	for _,v := range teams.Teams {
		fmt.Print(v)
		ih.Ps.SetPoints(strconv.Itoa(v),&flaghandler.Points{0,0,1700})
	}
	return "OK"
}
