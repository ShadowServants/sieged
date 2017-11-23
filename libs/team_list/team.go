package team_list

import (
	"errors"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type TeamList struct {
	Teams []int `json:"teams"`
}


type TeamIP struct {
	Id int `yaml:"id"`
	Name string `yaml:"name"`
	Server string `yaml:"server"`
	Network string `yaml:"network"`
}

type TeamsFile struct {
	Teams []TeamIP `yaml:"teams"`
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


func LoadTeamsList(in []byte) ([]TeamIP,error) {
	var team_file TeamsFile
	err := yaml.Unmarshal(in,&team_file)
	if err != nil {
		return nil,err
	}
	return team_file.Teams, nil

}