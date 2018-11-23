package status

import (
	"encoding/json"
	"errors"
	"sieged/internal/team"
)

func Loads(s string) (*team.Status, error) {
	tr := new(team.Status)
	if err := json.Unmarshal([]byte(s), tr); err != nil {
		return nil, errors.New("cant_loads_team_status_with_points")
	}
	return tr, nil
}
