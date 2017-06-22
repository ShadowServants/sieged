package main

import (
	"github.com/johnnovikov/hackforces/back/libs/storage"
	"encoding/json"
	"errors"
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
	st storage.Storage
}

func (ps *PointsStorage) GetPoints(key string) (*Points,error){
	s,err := ps.st.Get(key)
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
		ps.st.Set(key,data)
	}

}
