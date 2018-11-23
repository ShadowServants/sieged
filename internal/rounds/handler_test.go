package rounds

import (
	"encoding/json"
	"fmt"
	"sieged/internal/team"
	"testing"
)

func TestRoundHandler_TestTeam(t *testing.T) {
	pts := team.Score{Plus: 1, Minus: 1, Points: 1700}

	tr := team.Status{TeamId: 1, StatusMessage: "Checker failed", Status: "Down", Points: pts}
	responses := make([]team.Status, 0)
	responses = append(responses, tr)
	rr := Response{responses}
	byt, _ := json.Marshal(&rr)
	fmt.Print(string(byt))
}

