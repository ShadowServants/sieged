package score

import (
	"encoding/json"
	"errors"
	"sieged/internal/team"
	"sieged/pkg/storage"
)

func Loads(s string) (*team.Score, error) {
	p := team.Score{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil, errors.New("cant_loads_team_score")
	}
	return &p, nil
}

func Dumps(p *team.Score) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

type Storage struct {
	St storage.Storage
}

func (ps *Storage) GetPoints(key string) (*team.Score, error) {
	s, err := ps.St.Get(key)
	if err != nil {
		return nil, err
	}
	points, err := Loads(s)
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (ps *Storage) SetPoints(key string, score *team.Score) {
	if data, err := Dumps(score); err == nil {
		ps.St.Set(key, data)
	}
}
