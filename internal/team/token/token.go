package token

import (
	"encoding/json"
	"errors"
	"sieged/pkg/random"
)

type Token struct {
	Token  string `json:"token"`
	TeamId int    `json:"team_id"`
}

func (token *Token) Dump() (string, error) {
	byt, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

func New(teamId int) *Token {
	t := new(Token)
	t.TeamId = teamId
	t.Token = random.String(10,random.LowerCaseBytes + random.UpperCaseBytes + random.Digits)
	return t
}

func Loads(serialized string) (*Token, error) {
	p := new(Token)
	if err := json.Unmarshal([]byte(serialized), p); err != nil {
		return nil, errors.New("failed_deserialize_token")
	}
	return p, nil
}
