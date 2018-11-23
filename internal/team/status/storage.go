package status

import (
"fmt"
	"sieged/pkg/storage"
)

type Storage struct {
	st storage.Storage
}

func NewStorage(st storage.Storage) *Storage {
	t := new(Storage)
	t.st = st
	return t
}

func (r *Storage) GetStatus(teamId int, round int) string {
	keyString := fmt.Sprintf("%d:%d", teamId, round)
	a, _ := r.st.Get(keyString)
	return a
}

func (r *Storage) SetStatus(teamId, round int, status string) {
	keyString := fmt.Sprintf("%d:%d", teamId, round)
	r.st.Set(keyString, status)
}

