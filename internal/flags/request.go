package flags

import (
	"encoding/json"
	"errors"
)

type Request struct {
	Flag string `json:"flag"`
	Team int    `json:"team"`
}

func DumpRequest(tr *Request) (string, error) {
	byt, err := json.Marshal(tr)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

func LoadsRequest(s string) (*Request, error) {
	tr := Request{}
	if err := json.Unmarshal([]byte(s), &tr); err != nil {
		return nil, errors.New("json_unmarshal_error")
	}
	return &tr, nil
}
