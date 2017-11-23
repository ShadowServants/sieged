package flag_router

import (
	"strings"
	"errors"
	"fmt"
	"strconv"
	"net"
	"bufio"
	"net/http"
	"bytes"
	//"hackforces/libs/storage"
	"hackforces/libs/tcp_pool"
	"hackforces/libs/flagresponse"
	"hackforces/service_controller/flag_handler"
)


type FlagRouter struct {
	//IpStorage storage.Storage
	IpStorage map[string]string
	//Storage storage.Storage
	team_num int
	serviceMap map[string]string
	reason_map map[string]string
	VisualisationEnabled bool
	VisualisationUrl string
	pool map[string]*tcp_pool.TcpConnectionPool
}

func (fr *FlagRouter) SetVisualisation(url string) *FlagRouter {
	fr.VisualisationUrl = url
	fr.VisualisationEnabled = true
	return fr
}


func NewFlagRouter(team_num int) *FlagRouter{
	fr := new(FlagRouter)
	fr.team_num = team_num
	fr.reason_map =  map[string]string {
		"self":"That's your own flag",
		"invalid":"That's bad flag",
		"already_submitted":"You already submit this flag",
		"team_not_found":"We can`t find your team",
		"too_old":"This flag is too old",
		"not_ok":"Your service status is not good",
	}
	fr.VisualisationEnabled = false
	fr.serviceMap = make(map[string]string)
	fr.pool = make(map[string]*tcp_pool.TcpConnectionPool)
	fr.IpStorage = make(map[string]string)
	return fr

}

func (fr *FlagRouter) AddTeam(ip_cidr string,team_id string) {
		fr.IpStorage[ip_cidr] = team_id
}

//func (fr *FlagRouter) SetIpStorage(st storage.Storage) *FlagRouter{
//	fr.IpStorage = st
//	return fr
//}

func (fr *FlagRouter) getHandler(key string) (*tcp_pool.TcpConnectionPool,error) {
	if len(fr.serviceMap) != 0 {
		if pool, ok := fr.pool[key]; ok {
			return pool, nil
		}
		return nil, errors.New("Handler not found")
	}

	//data_str := fr.Storage.Get("handlers")
	//TODO::Add unjson str and return
	return nil,nil
}

func (fr *FlagRouter) getTeamNums() int {
	if fr.team_num != 0 { //Get from cache
		return fr.team_num
	}
	return 0
	//num, _ := fr.Storage.Get("team_num")
	//num_int , _ := strconv.Atoi(num)
	//fr.team_num = num_int
	//return num_int
}

func (fr *FlagRouter) RegisterHandler(service_prefix string,service_ip string) {
	fr.serviceMap[service_prefix] = service_ip
	ip_port := strings.Split(service_ip,":")
	ip, port := ip_port[0],ip_port[1]
	fr.pool[service_prefix] = tcp_pool.NewTcpPool(ip,port,fr.getTeamNums()*5)

}


func (fr *FlagRouter) CheckFlag(conn net.Conn,flag string, from int) string {
	fmt.Println("Start checking")
	tr := flaghandler.TeamRequest{Flag:flag, Team: from}
	data,err := flaghandler.DumpsTeamRequest(&tr)
	if err != nil {
		return "Bad data"
	}
	_ ,err = fmt.Fprintf(conn, data)
	if err != nil {
		fmt.Println("Service flag_handler is down ")
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
	if fr.VisualisationEnabled {
		fr.SendToVis(response)
	}
	fmt.Println(response)
	if sr.Successful{
		return fmt.Sprintf("Wow! You learned %d",sr.Delta)
	} else {
		return fmt.Sprintf("Sorry, but %s",fr.reason_map[sr.Reason])
	}

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

func (fr *FlagRouter) HandleRequest(flag string,ip string) string {
	from := fr.GetTeamIdByIp(ip)
	if from == -1 {
		return "Sorry, but we can`t find your team"
	}
	first_char := flag[0]
	handler, err := fr.getHandler(string(first_char))
	if err != nil {
		return "Handler not found"
	}

	conn , err := handler.Pool.Get()
	if err != nil {
		fmt.Println("Connection pool is full for service",first_char)
		return "Network is down"
	}
	return fr.CheckFlag(conn,flag,from)

}

func (fr *FlagRouter) GetTeamIdByIp(ip_raw string) int {
    ind := strings.Index(ip_raw,":")
    if ind <= 0{
      ind = 1
    }
    ip := ip_raw[:ind]

    ip_bytes := net.ParseIP(ip)
    fmt.Printf("From %s \n",ip)
    team_id_int := -1
    for key, value := range fr.IpStorage {
    	_, ipv4Net, err := net.ParseCIDR(key)
    	fmt.Printf("%s \n",ipv4Net.String())
    	if err != nil {
    		continue
		}
		if ipv4Net.Contains(ip_bytes) {
			v_int, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			team_id_int = v_int
		}

	}

	return team_id_int
}
