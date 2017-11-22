package flaghandler

import (
	"encoding/json"
	"errors"
	"hackforces/libs/storage"
)

func LoadsPoints(s string) (*Points, error) {
	p := Points{}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil,errors.New("Cant unmarshall json points")
    }
	return &p,nil
}

func DumpPoints(p *Points) (string, error) {
	byt, err := json.Marshal(p)
	if err != nil {
		return "",err
	}
	return string(byt),nil
}

type PointsStorage struct {
	St storage.Storage
}

func (ps *PointsStorage) GetPoints(key string) (*Points,error){
	s,err := ps.St.Get(key)
	if err != nil {
		return nil,err
	}
	points,err := LoadsPoints(s)
	if err != nil {
		return nil,err
	}
	return points,nil
}

func (ps *PointsStorage) SetPoints(key string,points *Points) {
	if data, err := DumpPoints(points); err == nil {
		ps.St.Set(key,data)
	}

}
