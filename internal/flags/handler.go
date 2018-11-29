package flags

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sieged/pkg/tcp"
)

var poolFullError = errors.New("connection_pool_full")
var poolReadError = errors.New("pool_read_error")

type Handler interface {
	CheckFlag(flag string, team int) (*Response, error)
}

type TcpHandler struct {
	pool *tcp.ConnectionPool
}

func (h *TcpHandler) CheckFlag(flag string, team int) (*Response, error) {
	conn, err := h.pool.Get()
	if err != nil {
		return nil, poolFullError
	}
	return h.checkFlag(conn, flag, team)

}

func (h *TcpHandler) checkFlag(conn net.Conn, flag string, team int) (*Response, error) {
	tr := Request{Flag: flag, Team: team}
	data, err := DumpRequest(&tr)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintf(conn, data)

	if err != nil {
		return nil, poolReadError
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, poolReadError
	}
	return LoadsResponse(response)
}
