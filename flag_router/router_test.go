package flag_router

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)


func buildTestFlagRouter() *FlagRouter{
	router_t := NewFlagRouter(2)
	//router_t.SetIpStorage(storage.NewSimpleStorage())
	router_t.IpStorage["127.0.1.1/24"] = "1"
	router_t.IpStorage["127.0.2.2/24"] = "2"
	return router_t
}

var router = buildTestFlagRouter()


func TestFlagRouter_GetTeamIdByIp(t *testing.T) {
	id := router.GetTeamIdByIp("127.0.1.1:5000")
	Convey("Checking that router check ip correctly",t,func() {
		So(id,ShouldEqual,1)
	})
}

func TestFlagRouter_GetTeamIdByIp_second(t *testing.T) {
	router.IpStorage = make(map[string]string)
	router.IpStorage["127.0.1.1/32"] = "1"
	router.IpStorage["127.0.2.2/32"] = "2"

	id := router.GetTeamIdByIp("127.0.1.1:5000")
	Convey("Checking that router check ip correctly",t,func() {
		So(id,ShouldEqual,1)
	})
}

