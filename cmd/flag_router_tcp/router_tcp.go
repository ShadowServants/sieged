package main

import (
	"fmt"
	"sieged/internal/flags"
	"sieged/pkg/helpers"
	"net"
	"os"
	"strings"
)

type TcpRouter struct {
	Fr *flags.Router
	Port string
	Host string
}

func (tr *TcpRouter) SetHost(host string) * TcpRouter {
	tr.Host = host
	return tr
}
func (tr *TcpRouter) SetPort(port string) * TcpRouter {
	tr.Port = port
	return tr
}

func (tr *TcpRouter) SetRouter(router *flags.Router) *TcpRouter {
	tr.Fr = router
	return tr
}

func (tr *TcpRouter) handleRequest(conn net.Conn) {
	conn.Write([]byte("Please enter flags, one flag per line \n"))
	buf := make([]byte, 200)
	n, err := conn.Read(buf)
	for err == nil {
		st := helpers.FromBytesToString(buf,n)
		ip := conn.RemoteAddr()
		//TODO: Flags missed ?
		st = strings.Split(st,"\n")[0]
		resp := tr.Fr.HandleRequest(st,ip.String())
		conn.Write([]byte(resp+"\n"))
		n, err = conn.Read(buf)
	}
	conn.Close()

}

func (tr *TcpRouter) StartPolling() {
	l, err := net.Listen("tcp",tr.Host+":"+tr.Port)
	if err != nil {
		helpers.FailOnError(err,"Tcp bind failed")
		os.Exit(1)
	}
	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		} else {
			go tr.handleRequest(conn)
		}
	}
}
