package client

import "reflect"

type HttpFacade struct {
	Method string
	Uri string
	Data interface{}
}

func (hb *HttpQueryBuilder) BuildUri() {
	return
}