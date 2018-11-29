package flags

import (
	"testing"
)

func TestLoadsResponse(t *testing.T) {
	r := StealResponse()
	r.SetInitiator(1).SetTarget(2).SetDelta(10.0)
}
