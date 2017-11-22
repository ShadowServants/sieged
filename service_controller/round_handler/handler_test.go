package round_handler

import (
	"testing"
	"github.com/jnovikov/hackforces/service_controllerndler"
	"encoding/json"
	"fmt"
)

func TestRoundHandler_TestTeam(t *testing.T) {
	pts := flaghandler.Points{Plus:1,Minus:1, Points: 1700}

	tr := TeamResponse{1,"Checker failed","Down",pts}
	responses := make([]TeamResponse,0)
	responses = append(responses,tr)
	rr := RoundResponse{responses}
	byt,_ := json.Marshal(&rr)
	fmt.Print(string(byt))


}
