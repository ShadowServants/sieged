package storage

import (
	"github.com/deckarep/golang-set"
)

type KeySet interface {
	Add(key string, value string)
	Check(key string, value string) bool
}

func NewSimpleKeySet() *SimpleKeySet {
	ks := new(SimpleKeySet)
	ks.Build()
	return ks
}

type SimpleKeySet struct {
	m map[string]mapset.Set
}

func (ks *SimpleKeySet) Build() {
	ks.m = make(map[string]mapset.Set)

}

func (ks *SimpleKeySet) Add(key string, value string )  {
	if ks.m[key] == nil {
		ks.m[key] = mapset.NewSet()
	}
	ks.m[key].Add(value)
}

func (ks *SimpleKeySet) Check(key string, value string) (bool) {
	if set, ok := ks.m[key]; ok {
		return set.Contains(value)
	} else {
		return false
	}
}


