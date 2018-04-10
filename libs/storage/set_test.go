package storage

import "testing"

func TestSimpleKeySet_Add(t *testing.T) {
	ks := SimpleKeySet{}
	ks.Build()
	ks.Add("lolkek","qwe")
	ks.Add("lolkek","kekus")
	if !ks.Check("lolkek","qwe") {
		t.Error("lolkek shoud contain qwe")
	}
	if !ks.Check("lolkek","kekus") {
		t.Error("lolkek shoud contain kekus")
	}
	if ks.Check("lolkek","kekus1") {
		t.Error("lolkek shoudnt contain kekus1")
	}
	if ks.Check("lolkek11","kekus") {
		t.Error("lolkek11 doesnt exist")
	}


}
