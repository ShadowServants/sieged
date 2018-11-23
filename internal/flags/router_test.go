package flags

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func buildTestFlagRouter() *Router {
	routerT := NewRouter(2)
	routerT.IpStorage["127.0.1.1/24"] = "1"
	routerT.IpStorage["127.0.2.2/24"] = "2"
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
