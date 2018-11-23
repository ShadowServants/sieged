package score

import (
	"sieged/internal/team"
	"sieged/pkg/storage"
	"testing"
)

func TestLoadsScore(t *testing.T) {
	d := `{"plus": 1,"minus": 1, "points":1700}`
	tr, err := Loads(d)
	if err != nil {
		t.Error(err.Error())
	}
	if tr.Plus != 1 || tr.Minus != 1 || tr.Points != 1700 {
		t.Error("Json undumps failed")
	}
}

func TestDumpScore(t *testing.T) {
	p := team.Score{Points: 1700}
	tr, err := Dumps(&p)
	if err != nil {
		t.Error(err.Error())
	}
	if tr != `{"plus":0,"minus":0,"points":1700}` {
		t.Error("Dump failed ", tr)
	}
}

func TestPointsStorage_GetPoints(t *testing.T) {
	simple := storage.SimpleStorage{}
	simple.Init()
	ps := Storage{&simple}
	p := team.Score{Points: 1700}
	ps.SetPoints("1", &p)
	d, err := ps.GetPoints("1")
	if err != nil {
		t.Error(err.Error())
	}
	if *d != p {
		t.Error("ScoreStorage storage didnt work")
	}
}

func TestPointsStorage_SetPoints(t *testing.T) {
	simple := storage.SimpleStorage{}
	simple.Init()
	ps := Storage{&simple}
	p := team.Score{Points: 1700}
	ps.SetPoints("1", &p)
	d, err := ps.GetPoints("1")
	if err != nil {
		t.Error(err.Error())
	}
	if *d != p {
		t.Error("ScoreStorage storage didnt work")
	}
}

