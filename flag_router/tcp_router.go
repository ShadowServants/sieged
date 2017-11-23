package flag_router

import (
	"net"
	"hackforces/libs/helpers"
	"os"
	"fmt"
)

type TcpRouter struct {
	Fr *FlagRouter
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

func (tr *TcpRouter) SetRouter(router *FlagRouter) *TcpRouter {
	tr.Fr = router
	return tr
}

func (tr *TcpRouter) handleRequest(conn net.Conn) {
	conn.Write([]byte("Please enter flags, one flag per line \n"))
	buf := make([]byte, 200)
    _, err := conn.Read(buf)
    for err == nil {
        st := helpers.FromBytesToString(buf)
        ip := conn.RemoteAddr()
		resp := tr.Fr.HandleRequest(st,ip.String())
		conn.Write([]byte(resp+"\n"))
        _, err = conn.Read(buf)
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
            //os.Exit(1)
        } else {
	        go tr.handleRequest(conn)

		}
    }
}
