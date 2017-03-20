package main

import (
    "fmt"
    "net"
    "os"
    "github.com/mediocregopher/radix.v2/pool"
        "../helpers"
)

var (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
    REDIS_POOL_SIZE = 10
    REDIS_PORT = "6379"
    REDIS_IP = "localhost"
)

type FlagHandler struct {
        listener net.Listener
	redis_pool *pool.Pool
}

func (fh *FlagHandler) Init(listener net.Listener, pool *pool.Pool)  {
        fh.listener = listener
        fh.redis_pool = pool
}

func (fh *FlagHandler) CheckFlag(flag string) string {
        conn, err := fh.redis_pool.Get()
        defer fh.redis_pool.Put(conn)
        if err != nil {
                fmt.Println("REDIS IS DOWN!!!!!!!!")
                return "Service is down"
        }

        result, err := conn.Cmd("GET", flag).Str()
        if err != nil {
                return "Flag not found"
        }
        return result



}

func (fh *FlagHandler) StartPolling() {
        for {
                conn, err := fh.listener.Accept()
                if err != nil {
                        panic(err)
                }
                go fh.handleRequest(conn)

        }
}

func (fh *FlagHandler) handleRequest(conn net.Conn) {
        buf := make([]byte, 50)
        _, err := conn.Read(buf)
        for err == nil {
                st := helpers.FromBytesToString(buf)
                conn.Write([]byte(fh.CheckFlag(st)))
                _, err = conn.Read(buf)
        }
        conn.Close()
}


func main() {

    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    p, err := pool.New("tcp", REDIS_IP+":"+REDIS_PORT, REDIS_POOL_SIZE)
    if err != nil {
        fmt.Println("Error acessing redis:", err.Error())
        os.Exit(1)
    }
    app := FlagHandler{l,p}
    app.StartPolling()

    // Close the listener when the application closes.
    defer app.listener.Close()
}

