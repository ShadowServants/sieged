package flaghandler

import (
	"testing"
	//"fmt"
	"github.com/jnovikov/hackforces/back/libs/storage"
)

func TestLoadsPoints(t *testing.T) {
	d := `{"plus": 1,"minus": 1, "points":1700}`
	tr,err := LoadsPoints(d)
	if err != nil {
		t.Error(err.Error())
	}
	if tr.Plus != 1 || tr.Minus != 1 || tr.Points != 1700  {
		t.Error("Json undumps failed")
	}
}

func TestDumpPoints(t *testing.T) {
	p := Points{0,0,1700}
	tr, err := DumpPoints(&p)
	if err != nil {
		t.Error(err.Error())
	}
	if tr != `{"plus":0,"minus":0,"points":1700}` {
		t.Error("Dump failed ",tr)
	}
}

func TestPointsStorage_GetPoints(t *testing.T) {
	simple := storage.SimpleStorage{}
	simple.Init()
	ps:= PointsStorage{&simple}
	p := Points{0,0,1700}
	ps.SetPoints("1",&p)
	d,err := ps.GetPoints("1")
	if err != nil {
		t.Error(err.Error())
	}
	if *d != p {
		t.Error("Points storage didnt work")
	}
}

func TestPointsStorage_SetPoints(t *testing.T) {
	simple := storage.SimpleStorage{}
	simple.Init()
	ps:= PointsStorage{&simple}
	p := Points{0,0,1700}
	ps.SetPoints("1",&p)
	d,err := ps.GetPoints("1")
	if err != nil {
		t.Error(err.Error())
	}
	if *d != p {
		t.Error("Points storage didnt work")
	}
}


