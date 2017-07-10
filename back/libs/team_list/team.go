package team_list

import (
	"errors"
	"encoding/json"
)

type TeamList struct {
	Teams []int `json:"teams"`
}



func LoadsTeamList(s string) (*TeamList, error) {
	p := TeamList{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil,errors.New("Cant unmarshall json points")
    }
	return &p,nil
}

func DumpTeamList(p *TeamList) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "",err
	}
	return string(byt),nil
}

