package flaghandler

import (
	"math"
	"sieged/internal/flags"
	"sieged/internal/team"
	"sieged/pkg/storage"
	"testing"
	//"fmt"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
)

func BuildTestFlagHandler() *FlagHandler {
	factory := NewFlagHandlerFactory()
	factory.SetFlagStorage(storage.NewSimpleStorage())
	factory.SetPointStorage(storage.NewSimpleStorage())
	factory.SetRoundStorage(storage.NewSimpleStorage())
	factory.SetStatusStorage(storage.NewSimpleStorage())
	factory.SetTeamFlagsSet(storage.NewSimpleKeySet())
	fHandler := factory.GetFlagHandler()
	fHandler.RoundSt.SetRound(1)
	fHandler.RoundDelta = 3
	fHandler.CurrentRound = 1
	fHandler.RoundCached = false
	fHandler.TeamNum = 3
	fHandler.Points.SetPoints("1", &team.Score{Points: 1700})
	fHandler.Points.SetPoints("2", &team.Score{Points: 1700})
	fHandler.Points.SetPoints("3", &team.Score{Points: 1700})

	return fHandler

}

var flagHandler = BuildTestFlagHandler()

func TestFlagHandler_calcDelta_smaller(t *testing.T) {
	delta := flagHandler.calcDelta(5, 10)
	Convey("Check delta if attacker pts < victim_pts", t, func() {
		So(delta, ShouldAlmostEqual, 4)
	})
}

func TestFlagHandler_calcDelta_almost_equal(t *testing.T) {
	flagHandler.TeamNum = 15
	delta := flagHandler.calcDelta(1200, 1150)
	Convey("Check delta if attacker pts < victim_pts", t, func() {
		So(math.Floor(delta), ShouldAlmostEqual, 11.0)
	})
}

func TestFlagHandler_calcDelta_2(t *testing.T) {
	flagHandler.TeamNum = 100
	delta := flagHandler.calcDelta(1200000, 1500)

	Convey("Check delta if attacker pts > victim pts", t, func() {
		So(math.Floor(delta), ShouldAlmostEqual, 3.0)
	})
	flagHandler.TeamNum = 3
}

func TestFlagHandler_calc(t *testing.T) {
	Convey("Check team 1 attacks team 2", t, func() {
		attacker := 1
		victim := 2
		res := flagHandler.calc(attacker, victim)
		attackerData, _ := flagHandler.GetTeamDataById(attacker)
		victimData, _ := flagHandler.GetTeamDataById(victim)

		So(math.Floor(res), ShouldAlmostEqual, 17)
		So(attackerData.Score.Plus, ShouldEqual, 1)
		So(attackerData.Score.Minus, ShouldEqual, 0)
		So(victimData.Score.Plus, ShouldEqual, 0)
		So(victimData.Score.Minus, ShouldEqual, 1)
		So(attackerData.Score.Points, ShouldAlmostEqual, 1718.9852583126342)
		So(victimData.Score.Points, ShouldAlmostEqual, 1682.7406742612416)
	})
	Convey("Check team 2 attacks team 1", t, func() {
		attacker := 2
		victim := 1
		res := flagHandler.calc(attacker, victim)
		So(math.Round(res), ShouldAlmostEqual, 17)
		So(flagHandler.Teams[attacker].Score.Plus, ShouldEqual, 1)
		So(flagHandler.Teams[attacker].Score.Minus, ShouldEqual, 1)
		So(flagHandler.Teams[victim].Score.Plus, ShouldEqual, 1)
		So(flagHandler.Teams[victim].Score.Minus, ShouldEqual, 1)
		So(flagHandler.Teams[attacker].Score.Points, ShouldAlmostEqual, 1701.9380468542802)
		So(flagHandler.Teams[victim].Score.Points, ShouldAlmostEqual, 1701.5331014098717)
	})
	Convey("Check points zero", t, func() {
		attacker := 1
		victim := 2
		flagHandler.Teams[attacker].Score.Points = 1
		flagHandler.Teams[victim].Score.Points = 1
		flagHandler.calc(attacker, victim)
		So(flagHandler.Teams[victim].Score.Points, ShouldEqual, 0)
	})

}

func TestFlagHandler_Build(t *testing.T) {
	flagHandler = BuildTestFlagHandler()
	td, _ := flagHandler.GetTeamDataById(1)
	Convey("Check that base points are ok", t, func() {
		So(td.Score.Points, ShouldEqual, 1700)
		So(td.Score.Minus, ShouldEqual, 0)
		So(td.Score.Plus, ShouldEqual, 0)
	})
	Convey("Check that base points are stored", t, func() {
		pts, er := flagHandler.Points.GetPoints("1")
		So(er, ShouldEqual, nil)
		So(pts.Points, ShouldEqual, 1700)
		So(pts.Plus, ShouldEqual, 0)
		So(pts.Minus, ShouldEqual, 0)
	})
}

func TestLoadsTeamRequest(t *testing.T) {
	d := `{"team": 1,"flag": "flagflag"}`
	tr, err := flags.LoadsRequest(d)
	Convey("Test json team requests loads correctly", t, func() {
		So(err, ShouldEqual, nil)
		So(tr.Flag, ShouldEqual, "flagflag")
		So(tr.Team, ShouldEqual, 1)
	})
	d = `{"team": "hkjkjk","flag": 1}`
	tr, err = flags.LoadsRequest(d)
	Convey("Check json team requests loads failed", t, func() {
		So(err, ShouldNotEqual, nil)
		So(tr, ShouldEqual, nil)

	})
}

func TestFlagHandler_SetCaptured(t *testing.T) {
	fl := BuildTestFlagHandler()
	tr2 := flags.Request{Flag: "WowSuchFlag", Team: 2}
	fl.SetCaptured(&tr2)
	Convey("Check that flag was set as captured", t, func() {
		So(fl.TeamFlagsSet.Check("2", "WowSuchFlag"), ShouldEqual, true)
	})
}

func TestFlagHandler_CheckFlag(t *testing.T) {
	fl := BuildTestFlagHandler()
	Convey("Check flag doenst exist", t, func() {
		So(fl.CheckFlag("Not flag"), ShouldEqual, nil)
	})
	Convey("Create flag and check it", t, func() {
		flag := flags.Data{Team: 2, Round: 2}
		fl.Flags.SetData("Really flag", &flag)
		So(fl.CheckFlag("Really flag").Team, ShouldEqual, 2)
	})
}

func TestFlagHandler_ValidateFlag(t *testing.T) {
	tr := flags.Request{Flag: "flagflag", Team: 1}
	flag := flags.Data{Team: 1, Round: 1}
	flagHandler.Flags.SetData("flagflag", &flag)
	Convey("Check your own flag", t, func() {
		ok, msg := flagHandler.ValidateFlag(&tr)
		So(ok, ShouldEqual, false)
		So(msg, ShouldEqual, SelfFLagMessage)
	})
	tr.Flag = "bad flag"
	Convey("Check bad flag", t, func() {
		ok, msg := flagHandler.ValidateFlag(&tr)
		So(ok, ShouldEqual, false)
		So(msg, ShouldEqual, BadFlagMessage)
	})
	tr.Flag = "Captured flag"
	flagHandler.SetCaptured(&tr)
	Convey("Check captured flag", t, func() {
		ok, msg := flagHandler.ValidateFlag(&tr)
		So(ok, ShouldEqual, false)
		So(msg, ShouldEqual, AlreadySubmitMessage)
	})
	tr.Team = 2
	tr.Flag = "flagflag"
	flagHandler.StatusStorage.SetStatus(2, 1, "Up")

	Convey("Check capture", t, func() {
		ok, msg := flagHandler.ValidateFlag(&tr)
		So(ok, ShouldEqual, true)
		So(msg, ShouldEqual, "")

	})
	tr.Team = 0
	tr.Flag = "flagflag"
	Convey("Check bad team", t, func() {
		ok, msg := flagHandler.ValidateFlag(&tr)
		So(ok, ShouldEqual, false)
		So(msg, ShouldEqual, TeamNotFoundMessage)
	})
}

func TestFlagHandler_HandleRequest(t *testing.T) {
	fl := BuildTestFlagHandler()
	fl.StatusStorage.SetStatus(2, 1, "Up")
	//Give first team 2 flags
	fd1 := flags.Data{Team: 1, Round: 1}
	fd2 := flags.Data{Team: 1, Round: 2}
	fl.Flags.SetData("flagteam1", &fd1)
	fl.Flags.SetData("flagteam1_1", &fd2)
	query := `{"team": 2,"flag": "notflag"}`
	Convey("Check team solve bad flag", t, func() {
		So(fl.HandleRequest(query), ShouldEqual, fmt.Sprintf(`{"successful":false,"type":"steal","initiator":2,"target":-1,"delta":0,"reason":"%s"}`, BadFlagMessage))
	})
	query = `{"team": 1,"flag": "flagteam1"}`
	Convey("Check team own flag", t, func() {
		So(fl.HandleRequest(query), ShouldEqual, fmt.Sprintf(`{"successful":false,"type":"steal","initiator":1,"target":-1,"delta":0,"reason":"%s"}`, SelfFLagMessage))

	})
	query = `{"team": 2,"flag": "flagteam1"}`
	Convey("Check captured flag", t, func() {

		So(fl.HandleRequest(query), ShouldEqual, `{"successful":true,"type":"steal","initiator":2,"target":1,"delta":17.259325738758434,"reason":""}`)
	})
	Convey("Check team try to catpure flag again", t, func() {
		So(fl.HandleRequest(query), ShouldEqual, fmt.Sprintf(`{"successful":false,"type":"steal","initiator":2,"target":-1,"delta":0,"reason":"%s"}`, AlreadySubmitMessage))

		//So(fl.HandleRequest(query),ShouldEqual,fmt.Sprintf(`{"ok":false,"text":"%s"}`,AlreadySubmitMessage))
		//So(fl.HandleRequest(query),ShouldEqual,AlreadySubmitMessage)
	})
	fl.RoundSt.SetRound(6)
	query = `{"team": 2,"flag": "flagteam1_1"}`
	Convey("Check flag is too old", t, func() {
		So(fl.HandleRequest(query), ShouldEqual, fmt.Sprintf(`{"successful":false,"type":"steal","initiator":2,"target":-1,"delta":0,"reason":"%s"}`, FlagTooOldMessage))

		//So(fl.HandleRequest(query),ShouldEqual,fmt.Sprintf(`{"ok":false,"text":"%s"}`,FlagTooOldMessage))
		//So(fl.HandleRequest(query),ShouldEqual,FlagTooOldMessage)
	})

}
