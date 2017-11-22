package flagstorage

import (
	"hackforces/libs/storage"
	"hackforces/libs/flagdata"
)

type FlagStorage struct {
	st storage.Storage
}

func NewFlagStorage(a storage.Storage) *FlagStorage {
	f := FlagStorage{}
	f.st = a
	return &f
}
func (ps *FlagStorage) GetFlagData(key string) (*flagdata.FlagData,error){
	s,err := ps.st.Get(key)
	if err != nil {
		return nil,err
	}
	flag,err := flagdata.LoadsFlagData(s)
	if err != nil {
		return nil,err
	}
	return flag,nil
}

func (ps *FlagStorage) SetFlagData(key string,flag *flagdata.FlagData) {
	if data, err := flagdata.DumpFlagData(flag); err == nil {
		ps.st.Set(key,data)
	}

}
