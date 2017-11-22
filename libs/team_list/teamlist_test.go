package team_list

import "testing"



func TestLoadsTeamList(t *testing.T) {
	d := `{"teams": [1, 2]}`
	tl, err := LoadsTeamList(d)
	if err != nil {
		t.Error(err.Error())
	}

	if tl.Teams[0] != 1 || tl.Teams[1] != 2{
		t.Error("Json undumps failed")

	}
}

func TestDumpPoints(t *testing.T) {
	tl := TeamList{[]int{0,1}}
	s, err := DumpTeamList(&tl)
	if err != nil {
		t.Error(err.Error())
	}
	if s != `{"teams":[0,1]}` {
		t.Error("Json dumps failed",s)
	}
}