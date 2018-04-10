package statusstorage

import (
	"fmt"
	"hackforces/libs/storage"
)

type StatusStorage struct {
	st storage.Storage
}

func NewStatusStorage(st storage.Storage) *StatusStorage{
	t := new(StatusStorage)
	t.st = st
	return t
}



func (r *StatusStorage) GetStatus(teamId int,round int) string {
	keyString := fmt.Sprintf("%d:%d", teamId,round)
	a, _ := r.st.Get(keyString)
	return a
}

func (r *StatusStorage) SetStatus(teamId, round int,status string) {
	keyString := fmt.Sprintf("%d:%d", teamId,round)
	r.st.Set(keyString,status)
}


