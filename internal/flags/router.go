package flags

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sieged/pkg/tcp"
	"strconv"
	"strings"
)

type Router struct {
	IpStorage            map[string]string
	teamNum              int
	reasonMap            map[string]string
	VisualisationEnabled bool
	VisualisationUrl     string
	handlerMap           map[string]Handler
	attacksLogger        *log.Logger
}

func (fr *Router) SetVisualisation(url string) *Router {
	fr.VisualisationUrl = url
	fr.VisualisationEnabled = true
	return fr
}

func (fr *Router) SetLogger(writer io.Writer) {
	fr.attacksLogger = log.New(writer, "Attack: ", log.Ldate|log.Ltime)
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
	fr.handlerMap = make(map[string]Handler)
	fr.IpStorage = make(map[string]string)
	return fr

}

func (fr *Router) AddTeam(ipCidr string, teamId string) {
	fr.IpStorage[ipCidr] = teamId
}

func (fr *Router) getHandler(key string) (Handler, error) {
	if len(fr.handlerMap) != 0 {
		if h, ok := fr.handlerMap[key]; ok {
			return h, nil
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

func (fr *Router) RegisterHandler(servicePrefix string, handler Handler) {
	fr.handlerMap[servicePrefix] = handler
}

func (fr *Router) RegisterTCPHandler(servicePrefix string, serviceIp string) {
	ipPort := strings.Split(serviceIp, ":")
	ip, port := ipPort[0], ipPort[1]
	handler := &TcpHandler{tcp.NewPool(ip, port, fr.getTeamNum()*5)}
	fr.RegisterHandler(servicePrefix, handler)
}

func (fr *Router) processSteal(attack *Response) {
	json, _ := DumpResponse(attack)
	if fr.VisualisationEnabled {
		fr.SendToVis(json)
	}
	if fr.attacksLogger == nil {
		fr.SetLogger(os.Stdout)
	}
	fr.attacksLogger.Println(json)
}

func (fr *Router) SendToVis(data string) {
	req, _ := http.NewRequest("POST", fr.VisualisationUrl, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send to visualisation")
	}

}

func (fr *Router) HandleFlag(flag string, id int) string {
	firstChar := flag[0]
	handler, err := fr.getHandler(string(firstChar))
	if err != nil {
		return "invalid flag"
	}

	response, err := handler.CheckFlag(flag, id)
	if err == poolFullError {
		log.Println("Connection pool is full for service", firstChar)
		return "try again later"
	}
	if err == poolReadError {
		log.Println("Connection is down for service", firstChar)
		return "try again later"
	}
	if err != nil {
		log.Println("Service:", firstChar, err.Error())
		return "try again later"
	}

	if response.Successful {
		fr.processSteal(response)
		return fmt.Sprintf("Accepted. %6f flag points", response.Delta)
	} else {
		return fmt.Sprintf("Denied: %s", fr.reasonMap[response.Reason])
	}
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
		log.Println("Failed to parse IP ", ipRaw)
		return -1
	}
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
