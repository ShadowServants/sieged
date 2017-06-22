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

func (st *SimpleStorage) Init() {
	st.data = make(map[string]string)
}

func (st *SimpleStorage) Get(key string) (string,error) {
	if data, ok := st.data[key]; ok {
		return data,nil
	} else {
		return "",errors.New("Key is missing")
	}
}

func (st *SimpleStorage) Set(key string,value string)  {
	st.data[key] = value
}

