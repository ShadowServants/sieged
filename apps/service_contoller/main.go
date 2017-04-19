package main

import (
	"net/http"
	"github.com/mediocregopher/radix.v2/pool"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"os/exec"
	"./flag_handler"
	"net"
)

var (
	TCP_LISTENER_HOST = "0.0.0.0"
	TCP_LISTENER_PORT = "31337"
    REDIS_POOL_SIZE = 70
    REDIS_PORT = "6379"
    REDIS_IP = "localhost"
	HANDLER_PORT = "8010"
	CONTROLLER_PORT = "8009"
)

type ServiceController struct {
        //listener net.Listener
	NumOfTeams int
	redis_pool *pool.Pool
	PORT string
}



func (sc *ServiceController) StartPolling() {
	http.HandleFunc("/round",sc.roundHandler)
	http.HandleFunc("/init",sc.initHandler)
	http.HandleFunc("/team_num",sc.numTeamHandler)
	//http.HandleFunc("/add_team",addHandler)
	http.ListenAndServe(":"+sc.PORT,nil)
}


func (sc *ServiceController) numTeamHandler(w http.ResponseWriter,r *http.Request){
	key := r.FormValue("key")
	if key != "a123a"{
		fmt.Fprint(w,"sorry")
		return
	}
	team_num_str := r.FormValue("team_num")
	team_num , _ :=strconv.Atoi(team_num_str)
	conn, _ := sc.redis_pool.Get()
	defer sc.redis_pool.Put(conn)
	conn.Cmd("auth","polinadrink")
	conn.Cmd("SET","teamnum",team_num)	
	sc.NumOfTeams = team_num
}

func (sc *ServiceController) initHandler(w http.ResponseWriter,r *http.Request) {
	conn, err := sc.redis_pool.Get()
	defer sc.redis_pool.Put(conn)
	conn.Cmd("auth","polinadrink")
	key := r.FormValue("key")
	if key != "a123a"{
		fmt.Fprint(w,"sorry")
		return
	}
	if err != nil {
		fmt.Println("REDIS IS DOWN!!!!!!!!")
		fmt.Fprint(w, "Service is down")
		return
	}

	team_num_str := r.FormValue("team_num")
	team_num , _ := strconv.Atoi(team_num_str)
	sc.NumOfTeams = team_num
	errors := make([]map[int]string,0)
	for i:=1 ; i <= team_num; i++ {
		_, err = conn.Cmd("HSET","team_points:plus",i,0).Int()
		_, err = conn.Cmd("HSET","team_points:minus",i,0).Int()
		_, err = conn.Cmd("HSET","team_points:fp",i,1700).Int()
		if err != nil {
			temp := map[int]string{i:string(err.Error())}
			errors = append(errors,temp)
		}
		_, err := conn.Cmd("HSET","statuses",i,"Down").Int()
		if err != nil{
			temp := map[int]string{i:string(err.Error())}
			errors = append(errors,temp)
		}
	}
	answer_json,_ := json.Marshal(errors)
	fmt.Fprint(w,string(answer_json))
	return



}

func (sc *ServiceController) getTeamIp(team_id string) string{
	conn, _ := sc.redis_pool.Get()
	defer sc.redis_pool.Put(conn)
	conn.Cmd("auth","polinadrink")
	result, _ := conn.Cmd("HGET","team_to_ip", team_id).Str()
	return result
}

func (sc *ServiceController) teamChecker(round_id,team_id string, output chan string) {
	ip := sc.getTeamIp(team_id)
	cmd := exec.Command("time_table.py","-t",team_id,"-r",round_id,"--ip",ip,"-tl","3")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		output <- "{\"status\": \"Down\", \"status_message\": \"\", \"team_id\": "+team_id+"}"
		return
	}
	fmt.Println(string(stdout))
	output <- string(stdout)
}

type CheckerResponse struct {
	team_id int
	status string
	status_message string
}


func (sc *ServiceController) roundHandler(w http.ResponseWriter,r *http.Request) {
	conn, _ := sc.redis_pool.Get()
	conn.Cmd("auth","polinadrink")
	defer sc.redis_pool.Put(conn)
	key := r.FormValue("key")
	if key != "a123a"{
		fmt.Fprint(w,"sorry")
		return
	}
	output := make(chan string,sc.NumOfTeams)
	round := r.FormValue("round")
	round_int, _ := strconv.Atoi(round)
	for i:=1; i <= sc.NumOfTeams ; i++{
		go sc.teamChecker(round,strconv.Itoa(i),output)
	}
	results := make([]string,0)
	for i:=1 ; i <= sc.NumOfTeams; i ++ {
		results = append(results,<-output)

	}
	response,_ := json.Marshal(results)
	fmt.Fprint(w,string(response))
	flag_handler.UpdateRound(round_int)
	

}

func main() {
	redis_pool,err  := pool.New("tcp", REDIS_IP+":"+REDIS_PORT, REDIS_POOL_SIZE)
	if err != nil {
		fmt.Println("Cant connect to redis",err.Error())
		os.Exit(1)
	}
	conn, _ := redis_pool.Get()
	defer redis_pool.Put(conn)
	conn.Cmd("auth","polinadrink")
	l, err := net.Listen("tcp", TCP_LISTENER_HOST+":"+TCP_LISTENER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	num,_ := conn.Cmd("get","teamnum").Int()
	fmt.Println("Controller started in",CONTROLLER_PORT)
	fmt.Println("Handler started in ",HANDLER_PORT)
	service := ServiceController{num,redis_pool,CONTROLLER_PORT}
	fl_handler := flag_handler.FlagHandler{Listener:l,Redis_pool:redis_pool,PORT:HANDLER_PORT}
	go fl_handler.StartPolling()
	service.StartPolling()



}
