package team

import (
	"encoding/json"
	"errors"
)

type List struct {
	Teams []int `json:"teams"`
}

type Address struct {
	Id      int    `yaml:"id"`
	Name    string `yaml:"name"`
	Server  string `yaml:"server"`
	Network string `yaml:"network"`
}

func LoadList(s string) (*List, error) {
	p := List{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil, errors.New("cant_loads_team_json")
	}
	return &p, nil
}

func DumpList(p *List) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

