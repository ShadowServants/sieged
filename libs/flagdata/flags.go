package flagdata

import (
	"encoding/json"
	"errors"
)

type FlagData struct {
	Team int `json:"team"`
	Round int `json:"round"`
}


func LoadsFlagData(s string) (*FlagData, error) {
	p := FlagData{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil,errors.New("Cant unmarshall json flagData")
    }
	return &p,nil
}

func DumpFlagData(p *FlagData) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "",err
	}
	return string(byt),nil
}