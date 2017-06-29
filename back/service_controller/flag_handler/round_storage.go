package main

import (
	"github.com/jnovikov/hackforces/back/libs/storage"
	"strconv"
)

type RoundStorage struct {
	st storage.Storage
}



func (r *RoundStorage) GetRound() int {
	a, _ := r.st.Get("round")
	round, _ := strconv.Atoi(a)
	return round
}

func (r *RoundStorage) SetRound(round int) {
	roundstr := strconv.Itoa(round)
	r.st.Set("round",roundstr)
}

