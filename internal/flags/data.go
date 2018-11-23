package flags

import (
	"encoding/json"
	"errors"
)

type Data struct {
	Team  int `json:"team"`
	Round int `json:"round"`
}

func LoadsData(s string) (*Data, error) {
	p := Data{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil, errors.New("failed_convert_json")
	}
	return &p, nil
}

func DumpData(p *Data) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}
