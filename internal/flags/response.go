package flags

import (
	"encoding/json"
	"errors"
)

type Response struct {
	Successful bool `json:"successful"`
	Type string `json:"type"`
	Initiator int `json:"initiator"`
	Target int `json:"target"`
	Delta float64 `json:"delta"`
	Reason string `json:"reason"`
}


func (hr *Response) SetType(typeStr string) *Response {
	hr.Type = typeStr
	return hr
}

func (hr *Response) SetInitiator(initiator int) *Response {
	hr.Initiator = initiator
	return hr
}

func (hr *Response) SetDelta(delta float64) *Response {
	hr.Delta = delta
	return hr
}

func (hr *Response) SetReason(reason string) *Response {
	hr.Reason = reason
	return hr
}

func (hr *Response) SetSuccessful(such bool) *Response {
	hr.Successful = such
	return hr
}

func (hr *Response) SetTarget(target int) *Response {
	hr.Target = target
	return hr
}

func StealResponse() *Response {
	hr := new(Response)
	hr.Type = "steal"
	return hr
}
func DumpResponse(p *Response) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "",err
	}
	return string(byt),nil
}

func LoadsResponse(s string) (*Response, error) {
	p := Response{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil,errors.New("cant_loads_flags_response")
	}
	return &p,nil
}
