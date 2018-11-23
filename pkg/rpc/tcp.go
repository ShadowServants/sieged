package rpc

import (
	"fmt"
	"sieged/pkg/helpers"
	"net"
	"os"
)

type TcpRpc struct {
	Port string
	Addr string
	Handler DataHandler
}


func (tr *TcpRpc) handleRequest(conn net.Conn) {
	buf := make([]byte, 200)
    n, err := conn.Read(buf)
    for err == nil {
        st := helpers.FromBytesToString(buf,n)
		resp := tr.Handler.HandleRequest(st)
		conn.Write([]byte(resp+"\n"))
        n, err = conn.Read(buf)
    }
    conn.Close()

}

func (tr *TcpRpc) Handle() {
	l, err := net.Listen("tcp",tr.Addr+":"+tr.Port)
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
