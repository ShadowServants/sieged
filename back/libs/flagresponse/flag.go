package flagresponse

import (
	"encoding/json"
	"errors"
)

type HandlerResponse struct {
	Successful bool `json:"successful"`
	Type string `json:"type"`
	Initiator int `json:"initiator"`
	Target int `json:"target"`
	Delta int `json:"delta"`
	Reason string `json:"reason"`
}


func DumpHandlerResponse(p *HandlerResponse) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "",err
	}
	return string(byt),nil
}

func LoadsHandlerResponse(s string) (*HandlerResponse, error) {
	p := HandlerResponse{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil,errors.New("Cant unmarshall json response")
    }
	return &p,nil
}
