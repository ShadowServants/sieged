package main

import (
	"testing"
	"github.com/jnovikov/hackforces/back/libs/storage"
	//"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/jnovikov/hackforces/back/libs/flagstorage"
	"github.com/jnovikov/hackforces/back/libs/flagdata"
)




func BuildTestFlagHandler() *FlagHandler {
	ps := storage.SimpleStorage{}
	ps.Init()
	pointStorage := PointsStorage{&ps}
	fs := storage.SimpleStorage{}
	fs.Init()

	flags := flagstorage.NewFlagStorage(&fs)
	ks := storage.SimpleKeySet{}
	ks.Build()
	fl := FlagHandler{}
	fl.TeamFlagsSet = &ks
	fl.Points = &pointStorage
	rs := storage.SimpleStorage{}
	rs.Init()
	roundst := &RoundStorage{&rs}
	fl.RoundSt = roundst
	fl.RoundSt.SetRound(1)
	fl.Flags = flags
	fl.RoundDelta = 3
	fl.CurrentRound = 1
	fl.RoundCached = false
	fl.Build(3,1700)
	return &fl
}

func TestFlagHandler_calcDelta(t *testing.T) {
	fl := BuildTestFlagHandler()
	delta := fl.calcDelta(1,1)
	Convey("Check delta",t,func(){
		So(delta,ShouldAlmostEqual,1)
	})
}



func TestFlagHandler_calc(t *testing.T) {
	fl := BuildTestFlagHandler()
	Convey("Check team 1 attacks team 2",t,func(){
		attacker := 1
		victim := 2
		res := fl.calc(attacker,victim)
		So(res,ShouldAlmostEqual,15)
		So(fl.Teams[attacker].points.Plus,ShouldEqual,1)
		So(fl.Teams[attacker].points.Minus,ShouldEqual,0)
		So(fl.Teams[victim].points.Plus,ShouldEqual,0)
		So(fl.Teams[victim].points.Minus,ShouldEqual,1)
		So(fl.Teams[attacker].points.Points,ShouldAlmostEqual,1715)
		So(fl.Teams[victim].points.Points,ShouldAlmostEqual,1685)
	})
	Convey("Check team 2 attacks team 1",t,func(){
		attacker := 2
		victim := 1
		res := fl.calc(attacker,victim)
		So(res,ShouldAlmostEqual,15)
		So(fl.Teams[attacker].points.Plus,ShouldEqual,1)
		So(fl.Teams[attacker].points.Minus,ShouldEqual,1)
		So(fl.Teams[victim].points.Plus,ShouldEqual,1)
		So(fl.Teams[victim].points.Minus,ShouldEqual,1)
		So(fl.Teams[attacker].points.Points,ShouldAlmostEqual,1700)
		So(fl.Teams[victim].points.Points,ShouldAlmostEqual,1700)
	})
	Convey("Check points zero",t,func(){
		attacker := 1
		victim := 2
		fl.Teams[attacker].points.Points = 1
		fl.Teams[victim].points.Points = 1
		fl.calc(attacker,victim)
		So(fl.Teams[victim].points.Points,ShouldEqual,0)
	})

}

func TestFlagHandler_Build(t *testing.T) {
	fl := BuildTestFlagHandler()
	td := fl.Teams[1]
	Convey("Check that base points are ok",t,func(){
		So(td.points.Points,ShouldEqual,1700)
		So(td.points.Minus,ShouldEqual,0)
		So(td.points.Plus,ShouldEqual,0)
	})
	Convey("Check that base points are stored",t,func(){
		pts,er := fl.Points.GetPoints("1")
		So(er,ShouldEqual,nil)
		So(pts.Points,ShouldEqual,1700)
		So(pts.Plus,ShouldEqual,0)
		So(pts.Minus,ShouldEqual,0)
	})
}

func TestLoadsTeamRequest(t *testing.T) {
	d := `{"team": 1,"flag": "flagflag"}`
	tr,err := LoadsTeamRequest(d)
	Convey("Test json team requests loads correctly",t,func(){
		So(err,ShouldEqual,nil)
		So(tr.Flag,ShouldEqual,"flagflag")
		So(tr.Team,ShouldEqual,1)
	})
	d = `{"team": "hkjkjk","flag": "flagflag","bad":"Bad"}`
	tr,err = LoadsTeamRequest(d)
	Convey("Check json team requests loads failed",t,func(){
		So(err,ShouldNotEqual,nil)
		So(tr,ShouldEqual,nil)

	})
}


func TestFlagHandler_SetCaptured(t *testing.T) {
	fl := BuildTestFlagHandler()
	tr2 := TeamRequest{"WowSuchFlag",2}
	fl.SetCaptured(&tr2)
	Convey("Check that flag was set as captured",t, func() {
		So(fl.TeamFlagsSet.Check("2","WowSuchFlag"),ShouldEqual,true)
	})
}


func TestFlagHandler_CheckFlag(t *testing.T) {
	fl := BuildTestFlagHandler()
	Convey("Check flag doenst exist",t,func () {
		So(fl.CheckFlag("Not flag"),ShouldEqual,nil)
	})
	Convey("Create flag and check it",t, func() {
		flag := flagdata.FlagData{2,2}
		fl.Flags.SetFlagData("Really flag",&flag)
		So(fl.CheckFlag("Really flag").Team,ShouldEqual,2)
	})
}

func TestFlagHandler_ValidateFlag(t *testing.T) {
	fl := BuildTestFlagHandler()
	tr := TeamRequest{"flagflag",1}
	flag := flagdata.FlagData{1,1}
	fl.Flags.SetFlagData("flagflag",&flag)
	Convey("Check your own flag",t,func(){
		ok, msg := fl.ValidateFlag(&tr)
		So(ok,ShouldEqual,false)
		So(msg,ShouldEqual,SelfFLagMessage)
	})
	tr.Flag = "bad flag"
	Convey("Check bad flag",t,func() {
		ok, msg := fl.ValidateFlag(&tr)
		So(ok,ShouldEqual,false)
		So(msg,ShouldEqual,BadFlagMessage)
	})
	tr.Flag = "Captured flag"
	fl.SetCaptured(&tr)
	Convey("Check captured flag",t,func() {
		ok, msg := fl.ValidateFlag(&tr)
		So(ok,ShouldEqual,false)
		So(msg,ShouldEqual,AlreadySubmitMessage)
	})
	tr.Team = 2
	tr.Flag = "flagflag"
	Convey("Check capture",t,func(){
		ok, msg := fl.ValidateFlag(&tr)
		So(ok,ShouldEqual,true)
		So(msg,ShouldEqual,"")

	})
	tr.Team = 0
	tr.Flag = "flagflag"
	Convey("Check bad team",t,func(){
		ok, msg := fl.ValidateFlag(&tr)
		So(ok,ShouldEqual,false)
		So(msg,ShouldEqual,TeamNotFoundMessage)
	})
}

func TestFlagHandler_HandleRequest(t *testing.T) {
	fl := BuildTestFlagHandler()
	Convey("Create simple flag handler with 3 teams", t, func() {
		So(len(fl.Teams),ShouldEqual,3)
	})
	//Give first team 2 flags
	fd1 := flagdata.FlagData{1,1}
	fd2 := flagdata.FlagData{1,2}
	fl.Flags.SetFlagData("flagteam1",&fd1)
	fl.Flags.SetFlagData("flagteam1_1",&fd2)
	query := `{"team": 2,"flag": "notflag"}`
	Convey("Check team solve bad flag",t,func(){
		So(fl.HandleRequest(query),ShouldEqual,BadFlagMessage)
	})
	query = `{"team": 1,"flag": "flagteam1"}`
	Convey("Check team own flag",t,func(){
		So(fl.HandleRequest(query),ShouldEqual,SelfFLagMessage)
	})
	query = `{"team": 2,"flag": "flagteam1"}`
	Convey("Check captured flag",t,func(){
		So(fl.HandleRequest(query),ShouldEqual,"Congrats. You captured 15 points")
	})
	Convey("Check team try to catpure flag again",t,func(){
		So(fl.HandleRequest(query),ShouldEqual,AlreadySubmitMessage)
	})
	fl.RoundSt.SetRound(6)
	query = `{"team": 2,"flag": "flagteam1_1"}`
	Convey("Check flag is too old",t,func(){
		So(fl.HandleRequest(query),ShouldEqual,FlagTooOldMessage)
	})


}

