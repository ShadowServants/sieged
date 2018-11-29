package flags

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

type MockHandler struct {
}

func (f *MockHandler) CheckFlag(flag string, team int) (*Response, error) {
	return StealResponse().SetSuccessful(true).SetDelta(1).SetInitiator(team).SetTarget(42), nil

}

func buildTestFlagRouter() *Router {
	routerT := NewRouter(2)
	routerT.IpStorage["127.0.1.1/24"] = "1"
	routerT.IpStorage["127.0.2.2/24"] = "2"
	routerT.RegisterHandler("T", new(MockHandler))
	return routerT
}

var router = buildTestFlagRouter()

func TestFlagRouter_GetTeamIdByIp(t *testing.T) {
	id := router.GetTeamIdByIp("127.0.1.1:5000")
	Convey("Checking that router check ip correctly", t, func() {
		So(id, ShouldEqual, 1)
	})
}

func TestFlagRouter_GetTeamIdByIp_second(t *testing.T) {
	router.IpStorage = make(map[string]string)
	router.IpStorage["127.0.1.1/32"] = "1"
	router.IpStorage["127.0.2.2/32"] = "2"

	id := router.GetTeamIdByIp("127.0.1.1:5000")
	Convey("Checking that router check ip correctly", t, func() {
		So(id, ShouldEqual, 1)
	})
}

func Test_Log(t *testing.T) {
	var b strings.Builder
	router.SetLogger(&b)
	router.HandleRequest("TESTTESTTEST", "127.0.1.1:5000")
	Convey("Checking log", t, func() {
		So(strings.Contains(b.String(), "Attack"), ShouldEqual, true)
	})

}
