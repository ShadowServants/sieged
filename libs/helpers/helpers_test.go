package helpers

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
)

func TestFromBytesToString(t *testing.T) {
	byt := []byte{0,1,2,3,97}
	FromBytesToString(byt,len(byt))
	convey.Convey("Assert that shitty bytes dont raise panic",t,func() {
		convey.So(true,convey.ShouldEqual,true)
	})
}
