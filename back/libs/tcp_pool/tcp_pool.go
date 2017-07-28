package tcp_pool

import (
	"net"
	"fmt"
	"gopkg.in/fatih/pool.v2"
	"github.com/jnovikov/hackforces/back/libs/helpers"
)

type TcpConnectionPool struct {
	factory func() (net.Conn, error)
	Pool pool.Pool
}

func NewTcpPool(host string,port string,size int) *TcpConnectionPool {
	t := new(TcpConnectionPool)
	addr := fmt.Sprintf("%s:%s",host,port)
	t.factory =  func() (net.Conn, error) { return net.Dial("tcp", addr) }
	p, err := pool.NewChannelPool(5, size, t.factory)
	helpers.FailOnError(err,"Cant create tcp pool")
	t.Pool = p
	return t

}
