package storage

import (
	"errors"
)

type Storage interface {
	Get(key string) (string,error)
	Set(key string, value string)
}



type SimpleStorage struct {
	data map[string] string
}

func (st *SimpleStorage) Init() *SimpleStorage{
	st.data = make(map[string]string)
	return st
}

func (st *SimpleStorage) Get(key string) (string,error) {
	if data, ok := st.data[key]; ok {
		return data,nil
	} else {
		return "",errors.New("key_is_missing")
	}
}

func (st *SimpleStorage) Set(key string,value string)  {
	st.data[key] = value
}


func NewSimpleStorage() *SimpleStorage{
	return new(SimpleStorage).Init()
}
