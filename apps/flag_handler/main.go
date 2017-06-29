package main

import (
    "fmt"
    "net"
    "github.com/mediocregopher/radix.v2/pool"
    "math"
    "encoding/json"
    "strings"
    "github.com/jnovikov/hackforces/apps/helpers"
    "strconv"
    "os"
    "net/http"
)

var (
    CONN_HOST = "0.0.0.0"
    CONN_PORT = "3333"
    HTTP_PORT = "3737"
    CONN_TYPE = "tcp"
    REDIS_POOL_SIZE = 30
    REDIS_PORT = "6378"
    REDIS_IP = "localhost"
)


type FlagHandler struct {
    Listener net.Listener
	Redis_pool *pool.Pool
}


func (fh *FlagHandler) GetTeamIdByIp(ip_raw string) int {
    fmt.Println(ip_raw)
    ind := strings.Index(ip_raw,":")
    if ind <= 0{
      ind = 1
    }
    ip := ip_raw[:ind]
    conn, _ := fh.Redis_pool.Get()
    defer fh.Redis_pool.Put(conn)
    team_id,_ := conn.Cmd("HGET","ip_to_team",ip).Int()
    fmt.Println(team_id)
    return team_id
}


func (fh *FlagHandler) Init(listener net.Listener, pool *pool.Pool)  {
        fh.Redis_pool = pool
}

func (fh *FlagHandler) GetRound() int {
    conn, _ := fh.Redis_pool.Get()
    defer fh.Redis_pool.Put(conn)
    result, err := conn.Cmd("GET","round_num").Int()
    fmt.Println(err)
    if err != nil {
        panic("Redis is down")
    }
    return result

}
func (fh *FlagHandler) CheckFlag(flag string,from int) string {
        conn, err := fh.Redis_pool.Get()
        defer fh.Redis_pool.Put(conn)
        if err != nil {
                fmt.Println("REDIS IS DOWN!!!!!!!!")
                return "Service is down"
        }

        result, err := conn.Cmd("HGET","flags", flag).Str()
        if err != nil {
                return "Flag not found \n"
        }
        status, err := conn.Cmd("HGET","statuses",from).Str()
        if err != nil {
            return "Status not founded :("
        }
        dict:= make(map[string]int)
        json.Unmarshal([]byte(result),&dict)
        victim_id := dict["team_id"]
        if victim_id == from{
            return "Your own flag :( \n"
        }
        fmt.Println(victim_id,from)
        fmt.Println(status)
        if status != "Up" {
           return "Your service status is not Ok\n"
        }

        ownded,_ := conn.Cmd("HGET","owned:"+strconv.Itoa(from),flag).Int()
        if ownded == 1{
            return "You already solved this flag"
        }
        round_id := fh.GetRound()
        if round_id - dict["round_id"] > 3 {
            return "Flag is too old\n"
        }
        conn.Cmd("HINCRBY","team_points:plus",from,1)
        conn.Cmd("HINCRBY","team_points:minus", victim_id,-1)
        conn.Cmd("HSET","owned:"+strconv.Itoa(from),flag,1)
        attacker_points,_ := conn.Cmd("HGET","team_points:fp",from).Int()
        victim_points, _ := conn.Cmd("HGET","team_points:fp",victim_id).Int()
        ap := math.Max(1.0,float64(attacker_points + 1))
        vp := math.Max(1.0,float64(victim_points + 1))
        logattacker := math.Log2(ap) + 1
        logvictim := math.Log2(vp) + 1
        delta := logvictim / logattacker
        delta_points := int(delta * 15)
        conn.Cmd("HINCRBY","team_points:fp",from,delta_points).Int()
        conn.Cmd("HINCRBY","team_points:fp",victim_id,-1*delta_points).Int()
        return "Congrats. You learned" + strconv.Itoa(delta_points) +"\n"

}

func (fh *FlagHandler) HandleTCPRequest(conn net.Conn) {
    buf := make([]byte, 50)
    _, err := conn.Read(buf)
    for err == nil {
        st := helpers.FromBytesToString(buf)
        resp := fh.CheckFlag(st,fh.GetTeamIdByIp(conn.RemoteAddr().String()))
        conn.Write([]byte(resp))
        _, err = conn.Read(buf)
    }
    conn.Close()
}

func (fh *FlagHandler) StartTCPListen() {
    for {
        conn, err := fh.Listener.Accept()
        if err != nil {
            fmt.Println(err.Error())
        }
        go fh.HandleTCPRequest(conn)
    }
}

//func (fh *FlagHandler) StartPolling() {
//        fh.StartTCPListen()
//}


func (fh *FlagHandler) StartPolling() {
        http.HandleFunc("/",fh.handleRequest)
        go fh.StartTCPListen()
        http.ListenAndServe("0.0.0.0:"+HTTP_PORT,nil)
}

func (fh *FlagHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
        flag := r.FormValue("flag")
        ip := r.RemoteAddr
        from := fh.GetTeamIdByIp(ip)
        fmt.Fprint(w,fh.CheckFlag(flag,from))
}


func main() {

    p, err := pool.New("tcp", REDIS_IP+":"+REDIS_PORT, REDIS_POOL_SIZE)
    if err != nil {
        fmt.Println("Error acessing redis:", err.Error())
        os.Exit(1)
    }
    l, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
    app := FlagHandler{l,p,}

    app.StartPolling()

    // Close the listener when the application closes.
    //defer app.listener.Close()
}

