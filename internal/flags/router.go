package flags

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sieged/pkg/tcp"
	"strconv"
	"strings"
)

type Router struct {
	IpStorage            map[string]string
	teamNum              int
	serviceMap           map[string]string
	reasonMap            map[string]string
	VisualisationEnabled bool
	VisualisationUrl     string
	pool                 map[string]*tcp.ConnectionPool
}

func (fr *Router) SetVisualisation(url string) *Router {
	fr.VisualisationUrl = url
	fr.VisualisationEnabled = true
	return fr
}

func NewRouter(teamNum int) *Router {
	fr := new(Router)
	fr.teamNum = teamNum
	fr.reasonMap = map[string]string{
		"self":              "flags is your own",
		"invalid":           "invalid flag",
		"already_submitted": "flag is already submitted",
		"team_not_found":    "team not found",
		"too_old":           "flag is too old",
	}
	fr.VisualisationEnabled = false
	fr.serviceMap = make(map[string]string)
	fr.pool = make(map[string]*tcp.ConnectionPool)
	fr.IpStorage = make(map[string]string)
	return fr

}

func (fr *Router) AddTeam(ipCidr string, teamId string) {
	fr.IpStorage[ipCidr] = teamId
}

func (fr *Router) getHandler(key string) (*tcp.ConnectionPool, error) {
	if len(fr.serviceMap) != 0 {
		if pool, ok := fr.pool[key]; ok {
			return pool, nil
		}
		return nil, errors.New("handler not found")
	}
	return nil, errors.New("no services found")
}

func (fr *Router) getTeamNum() int {
	if fr.teamNum != 0 {
		return fr.teamNum
	}
	return 0
}

func (fr *Router) RegisterHandler(servicePrefix string, serviceIp string) {
	fr.serviceMap[servicePrefix] = serviceIp
	ipPort := strings.Split(serviceIp, ":")
	ip, port := ipPort[0], ipPort[1]
	fr.pool[servicePrefix] = tcp.NewPool(ip, port, fr.getTeamNum()*5)

}

func (fr *Router) CheckFlag(conn net.Conn, flag string, from int) string {
	tr := Request{Flag: flag, Team: from}
	data, err := DumpRequest(&tr)
	if err != nil {
		return "invalid data format"
	}
	_, err = fmt.Fprintf(conn, data)
	if err != nil {
		fmt.Println("Service flag_handler is down ")
		return "try again later"
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Cant read response from flag_handler")
		return "try again later"
	}
	sr, err := LoadsResponse(response)
	if err != nil {
		fmt.Println("Cant loads response from flag_handler")
		return "try again later"
	}
	if fr.VisualisationEnabled {
		fr.SendToVis(response)
	}
	if sr.Successful {
		return fmt.Sprintf("Accepted. %6f flag points", sr.Delta)
	} else {
		return fmt.Sprintf("Denied: %s", fr.reasonMap[sr.Reason])
	}

}

func (fr *Router) SendToVis(data string) {
	req, _ := http.NewRequest("POST", fr.VisualisationUrl, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send to visualisation")
	}

}

func (fr *Router) HandleFlag(flag string, id int) string {
	firstChar := flag[0]
	handler, err := fr.getHandler(string(firstChar))
	if err != nil {
		return "invalid flag"
	}

	conn, err := handler.Pool.Get()
	if err != nil {
		fmt.Println("Connection pool is full for service", firstChar)
		return "network is down try again later"
	}
	return fr.CheckFlag(conn, flag, id)
}

func (fr *Router) HandleRequest(flag string, ip string) string {
	if len(flag) <= 0 {
		return "invalid data"
	}
	from := fr.GetTeamIdByIp(ip)
	if from == -1 {
		return "your IP is unknown"
	}
	return fr.HandleFlag(flag, from)
}

func (fr *Router) GetTeamIdByIp(ipRaw string) int {
	ind := strings.Index(ipRaw, ":")
	if ind <= 0 {
		ind = 1
	}
	ip := ipRaw[:ind]
	ipBytes := net.ParseIP(ip)
	if ipBytes == nil {
		fmt.Println("Failed to parse IP ", ipRaw)
		return -1
	}
	fmt.Printf("From %s \n", ip)
	teamIdInt := -1
	for key, value := range fr.IpStorage {
		_, ipv4Net, err := net.ParseCIDR(key)
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
