package token

import (
	"hackforces/libs/storage"
	"hackforces/libs/helpers"
	"strconv"
)

type Storage struct {
	st storage.Storage
}

func (ts *Storage) Find(key string) (*Token, error) {
	s, e := ts.st.Get(key)
	if e != nil {
		return nil, e
	}
	t, e := Loads(s)
	if e != nil {
		return nil, e
	}
	return t, e
}

func (ts *Storage) FindById(id int) (*Token, error) {
	key := "reverse_" + strconv.Itoa(id)
	s, e := ts.st.Get(key)
	if e != nil {
		return nil, e
	}
	t, e := Loads(s)
	if e != nil {
		return nil, e
	}
	return t, e
}

func (ts *Storage) New(teamId int) *Token {
	token := New(teamId)
	s, err := token.Dump()
	helpers.FailOnError(err, "Failed to dump token")
	ts.st.Set(token.Token, s)
	ts.st.Set("reverse_"+strconv.Itoa(teamId), s)
	return token
}

func NewStorage(st storage.Storage) *Storage {
	return &Storage{st}
}
