package statusstorage

import (
	"github.com/jnovikov/hackforces/libs/storage"
	"fmt"
)

type StatusStorage struct {
	st storage.Storage
}

func NewStatusStorage(st storage.Storage) *StatusStorage{
	t := new(StatusStorage)
	t.st = st
	return t
}



func (r *StatusStorage) GetStatus(team_id int,round int) string {
	key_string := fmt.Sprintf("%d:%d",team_id,round)
	a, _ := r.st.Get(key_string)
	return a
}

func (r *StatusStorage) SetStatus(team_id, round int,status string) {
	key_string := fmt.Sprintf("%d:%d",team_id,round)
	r.st.Set(key_string,status)
}


