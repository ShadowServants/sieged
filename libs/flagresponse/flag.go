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
	Delta float64 `json:"delta"`
	Reason string `json:"reason"`
}


func (hr *HandlerResponse) SetType(typeStr string) *HandlerResponse {
	hr.Type = typeStr
	return hr
}

func (hr *HandlerResponse) SetInitiator(initiator int) *HandlerResponse {
	hr.Initiator = initiator
	return hr
}

func (hr *HandlerResponse) SetDelta(delta float64) *HandlerResponse {
	hr.Delta = delta
	return hr
}

func (hr *HandlerResponse) SetReason(reason string) *HandlerResponse {
	hr.Reason = reason
	return hr
}

func (hr *HandlerResponse) SetSuccessful(such bool) *HandlerResponse {
	hr.Successful = such
	return hr
}

func (hr *HandlerResponse) SetTarget(target int) *HandlerResponse {
	hr.Target = target
	return hr
}

func NewStealResponse() *HandlerResponse {
	hr := new(HandlerResponse)
	hr.Type = "steal"
	return hr
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
