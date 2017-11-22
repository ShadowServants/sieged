package main

import (
	"github.com/jnovikov/hackforces/libs/storage"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"net"
	"gopkg.in/fatih/pool.v2"
	"github.com/jnovikov/hackforces/libs/helpers"
	"github.com/jnovikov/hackforces/service_controller/flag_handler"
	"bufio"
	"github.com/jnovikov/hackforces/libs/flagresponse"
	"bytes"
)

var reason_map = map[string]string {
	"self":"That's your own flag",
	"invalid":"That's bad flag",
	"already_submitted":"You already submit this flag",
	"team_not_found":"We can`t find your team",
	"too_old":"This flag is too old",
	"not_ok":"Your service status is not good",
}


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



type FlagRouter struct {
	Conn_pool *TcpConnectionPool
	Port string
	IpStorage storage.Storage
	VisualisationUrl string
	VisualisationEnabled bool
}

func (fh *FlagRouter) GetTeamIdByIp(ip_raw string) int {
    fmt.Println(ip_raw)
    ind := strings.Index(ip_raw,":")
    if ind <= 0{
      ind = 1
    }
    ip := ip_raw[:ind]
	team_id,err := fh.IpStorage.Get(ip)
	fmt.Println(ip,team_id)
	if err != nil {
		return -1
	}
	team_id_int,err := strconv.Atoi(team_id)
	if err != nil {
		return -1
	}
	return team_id_int
}

func (fh *FlagRouter) SendToVis(data string) {
	req, _ := http.NewRequest("POST", fh.VisualisationUrl, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil{
		fmt.Println("Failed to send to visualisation")
	}

}

func (fh *FlagRouter) CheckFlag(flag string, from int) string {
	tr := flaghandler.TeamRequest{flag,from}
	conn, err := fh.Conn_pool.Pool.Get()
	if err != nil {
		return "Network is down"
	}
	data,err := flaghandler.DumpsTeamRequest(&tr)
	if err != nil {
		return "Bad data"
	}
	_ ,err = fmt.Fprintf(conn, data)
	if err != nil {
		return "Network is down"
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "Network is down"
	}
	sr,err := flagresponse.LoadsHandlerResponse(response)
	if err != nil {
		return "Bad data"
	}
	if fh.VisualisationEnabled {
		fh.SendToVis(response)
	}
	fmt.Println(response)
	if sr.Successful{
		return fmt.Sprintf("Wow! You learned %d",sr.Delta)
	} else {
		return fmt.Sprintf("Sorry, but %s",reason_map[sr.Reason])
	}

}

func (fh *FlagRouter) handleRequest(w http.ResponseWriter, r *http.Request) {
	flag := r.FormValue("flag")
	ip := r.RemoteAddr
	from := fh.GetTeamIdByIp(ip)
	if from == -1 {
		fmt.Fprint(w,"Sorry, but we cant find your team")
		return
	}
	fmt.Fprint(w,fh.CheckFlag(flag,from))
}

func (fh *FlagRouter) StartPolling() {
		Rp := new(storage.RadixPool)
		Rp.Build("127.0.0.1","6379",27)
		fh.IpStorage = &storage.HsetRadixStorage{Rp,"player_ip_to_team"}
		fh.Conn_pool = NewTcpPool("127.0.0.1","7878",27)
        http.HandleFunc("/",fh.handleRequest)
        http.ListenAndServe("0.0.0.0:"+fh.Port,nil)
}


func main() {
	fr := new(FlagRouter)
	fr.VisualisationUrl = "http://127.0.0.1:3000/broadcast"
	fr.VisualisationEnabled = true
	fr.Port = "7331"
	fr.StartPolling()
}
