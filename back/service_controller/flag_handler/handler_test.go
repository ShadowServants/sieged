package main

import (
	"testing"
	//"fmt"
)

func TestLoadsTeamRequest(t *testing.T) {
	d := `{"team": 1,"flag": "flagflag"}`
	tr,err := LoadsTeamRequest(d)
	if err != nil {
		t.Error(err.Error())
	}
	if tr.Flag != "flagflag" || tr.Team != 1  {
		t.Error("Bad data")
	}
}



