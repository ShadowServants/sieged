package flags

import (
	"sieged/pkg/storage"
)

type Storage struct {
	st storage.Storage
}

func NewStorage(a storage.Storage) *Storage {
	f := Storage{}
	f.st = a
	return &f
}

func (ps *Storage) GetData(key string) (*Data, error){
	s, err := ps.st.Get(key)
	if err != nil {
		return nil, err
	}
	flag, err := LoadsData(s)
	if err != nil {
		return nil, err
	}
	return flag, nil
}

func (ps *Storage) SetData(key string, flag *Data) {
	if data, err := DumpData(flag); err == nil {
		ps.st.Set(key,data)
	}

}
