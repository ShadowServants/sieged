package tcp

import (
	"fmt"
	"gopkg.in/fatih/pool.v2"
	"sieged/pkg/helpers"
	"net"
)

type ConnectionPool struct {
	factory func() (net.Conn, error)
	Pool    pool.Pool
}

func NewPool(host string, port string, size int) *ConnectionPool {
	t := new(ConnectionPool)
	addr := fmt.Sprintf("%s:%s", host, port)
	t.factory = func() (net.Conn, error) { return net.Dial("tcp", addr) }
	p, err := pool.NewChannelPool(5, size, t.factory)
	helpers.FailOnError(err, "Cant create tcp pool")
	t.Pool = p
	return t
}
