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
	teamNum              int
	serviceMap           map[string]string
	reasonMap            map[string]string
	VisualisationEnabled bool
	VisualisationUrl     string
	pool                 map[string]*tcp_pool.TcpConnectionPool
}

func (fr *FlagRouter) SetVisualisation(url string) *FlagRouter {
	fr.VisualisationUrl = url
	fr.VisualisationEnabled = true
	return fr
}

func NewFlagRouter(teamNum int) *FlagRouter {
	fr := new(FlagRouter)
	fr.teamNum = teamNum
	fr.reasonMap = map[string]string{
		"self":              "That's your own flag",
		"invalid":           "That's bad flag",
		"already_submitted": "You already submit this flag",
		"team_not_found":    "We can`t find your team",
		"too_old":           "This flag is too old",
		"not_ok":            "Your service status is not good",
	}
	fr.VisualisationEnabled = false
	fr.serviceMap = make(map[string]string)
	fr.pool = make(map[string]*tcp_pool.TcpConnectionPool)
	fr.IpStorage = make(map[string]string)
	return fr

}

func (fr *FlagRouter) AddTeam(ipCidr string, teamId string) {
	fr.IpStorage[ipCidr] = teamId
}

//func (fr *FlagRouter) SetIpStorage(st storage.Storage) *FlagRouter{
//	fr.IpStorage = st
//	return fr
//}

func (fr *FlagRouter) getHandler(key string) (*tcp_pool.TcpConnectionPool, error) {
	if len(fr.serviceMap) != 0 {
		if pool, ok := fr.pool[key]; ok {
			return pool, nil
		}
		return nil, errors.New("handler not found")
	}

	//data_str := fr.Storage.Get("handlers")
	//TODO::Add unjson str and return
	return nil, errors.New("no services found")
}

func (fr *FlagRouter) getTeamNums() int {
	if fr.teamNum != 0 { //Get from cache
		return fr.teamNum
	}
	return 0
	//num, _ := fr.Storage.Get("team_num")
	//num_int , _ := strconv.Atoi(num)
	//fr.team_num = num_int
	//return num_int
}

func (fr *FlagRouter) RegisterHandler(servicePrefix string, serviceIp string) {
	fr.serviceMap[servicePrefix] = serviceIp
	ipPort := strings.Split(serviceIp, ":")
	ip, port := ipPort[0], ipPort[1]
	fr.pool[servicePrefix] = tcp_pool.NewTcpPool(ip, port, fr.getTeamNums()*5)

}

func (fr *FlagRouter) CheckFlag(conn net.Conn, flag string, from int) string {
	fmt.Println("Start checking")
	tr := flaghandler.TeamRequest{Flag: flag, Team: from}
	data, err := flaghandler.DumpsTeamRequest(&tr)
	if err != nil {
		return "Bad data"
	}
	_, err = fmt.Fprintf(conn, data)
	if err != nil {
		fmt.Println("Service flag_handler is down ")
		return "Network is down"
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "Network is down"
	}
	sr, err := flagresponse.LoadsHandlerResponse(response)
	if err != nil {
		return "Bad data"
	}
	if fr.VisualisationEnabled {
		fr.SendToVis(response)
	}
	fmt.Println(response)
	if sr.Successful {
		return fmt.Sprintf("Wow! You learned %6f", sr.Delta)
	} else {
		return fmt.Sprintf("Sorry, but %s", fr.reasonMap[sr.Reason])
	}

}

func (fr *FlagRouter) SendToVis(data string) {
	req, _ := http.NewRequest("POST", fr.VisualisationUrl, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send to visualisation")
	}

}

func (fr *FlagRouter) HandleFlag(flag string, id int) string {
	firstChar := flag[0]
	handler, err := fr.getHandler(string(firstChar))
	if err != nil {
		return "Invalid flag"
	}

	conn, err := handler.Pool.Get()
	if err != nil {
		fmt.Println("Connection pool is full for service", firstChar)
		return "Network is down"
	}
	return fr.CheckFlag(conn, flag, id)
}

func (fr *FlagRouter) HandleRequest(flag string, ip string) string {
	if len(flag) <= 0 {
		return "Bad"
	}
	from := fr.GetTeamIdByIp(ip)
	if from == -1 {
		return "Sorry, but we can`t find your team"
	}
	return fr.HandleFlag(flag, from)
}

func (fr *FlagRouter) GetTeamIdByIp(ipRaw string) int {
	ind := strings.Index(ipRaw, ":")
	if ind <= 0 {
		ind = 1
	}
	ip := ipRaw[:ind]

	ipBytes := net.ParseIP(ip)
	fmt.Printf("From %s \n", ip)
	teamIdInt := -1
	for key, value := range fr.IpStorage {
		_, ipv4Net, err := net.ParseCIDR(key)
		fmt.Printf("%s \n", ipv4Net.String())
		if err != nil {
			continue
		}
		if ipv4Net.Contains(ipBytes) {
			vInt, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			teamIdInt = vInt
		}

	}

	return teamIdInt
}
