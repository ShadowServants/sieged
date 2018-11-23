package rounds

import (
	"sieged/pkg/storage"
	"strconv"
)

type Storage struct {
	St storage.Storage
}

func (r *Storage) GetRound() int {
	a, _ := r.St.Get("round")
	round, _ := strconv.Atoi(a)
	return round
}

func (r *Storage) SetRound(round int) {

	roundStr := strconv.Itoa(round)
	r.St.Set("round", roundStr)
}
