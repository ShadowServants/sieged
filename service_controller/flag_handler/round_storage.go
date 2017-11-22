package flaghandler

import (
	"strconv"
	"hackforces/libs/storage"
)

type RoundStorage struct {
	St storage.Storage
}



func (r *RoundStorage) GetRound() int {
	a, _ := r.St.Get("round")
	round, _ := strconv.Atoi(a)
	return round
}

func (r *RoundStorage) SetRound(round int) {
	roundstr := strconv.Itoa(round)
	r.St.Set("round",roundstr)
}

