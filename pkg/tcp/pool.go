package tcp

import (
	"fmt"
	"gopkg.in/fatih/pool.v2"
	"net"
	"sieged/pkg/helpers"
)

type ConnectionPool struct {
	factory func() (net.Conn, error)
	pool.Pool
}

func NewPool(host string, port string, size int) *ConnectionPool {
	addr := fmt.Sprintf("%s:%s", host, port)
	factory := func() (net.Conn, error) { return net.Dial("tcp", addr) }
	p, err := pool.NewChannelPool(5, size, factory)
	helpers.FailOnError(err, "Cant create tcp pool")
	return &ConnectionPool{factory, p}
}
