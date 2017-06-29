package main

import (
	"net/http"
	"github.com/mediocregopher/radix.v2/pool"
	"fmt"
	"os"
	"math/rand"
	"time"
	"encoding/json"
	"strconv"
)
var (
    PORT = "8011"
    REDIS_POOL_SIZE = 30
    REDIS_PORT = "6378"
    REDIS_IP = "localhost"
)


var Redis_pool *pool.Pool
const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const serviceChar = "T"

func RandString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }
    return string(b)
}

func GenerateFlag() string {
	return serviceChar + RandString(30) + "="
}



func indexHandler(w http.ResponseWriter,r *http.Request){
	team_id := r.FormValue("team_id")
	round_id := r.FormValue("round_id")
	fmt.Println(team_id,round_id)
	if team_id == "" || round_id == ""{
		fmt.Fprint(w,"You shuld set round and team id.")
		return
	}
	team_id_int,_ := strconv.Atoi(team_id)
	round_id_int, _ := strconv.Atoi(round_id)
	flag_map := map[string]int{"team_id": team_id_int, "round_id": round_id_int}
	flag_json,_ := json.Marshal(flag_map)
	flag := GenerateFlag()
	conn, err := Redis_pool.Get()
	conn.Cmd("auth","polinadrink")
	if err != nil {
		fmt.Fprint(w,"Redis is down")
		return
	}
        defer Redis_pool.Put(conn)
	_, err = conn.Cmd("HSET","flags",flag,flag_json).Int()
	if err != nil{
		fmt.Println(err)
		fmt.Fprint(w,"Bad command")
		return
	}

	fmt.Fprint(w,flag)
}

func main() {
	var err error
	Redis_pool,err  = pool.New("tcp", REDIS_IP+":"+REDIS_PORT, REDIS_POOL_SIZE)
	if err != nil {
		fmt.Println("Cant connect to redis",err.Error())
		os.Exit(1)
	}
	fmt.Println("Started on",PORT)
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/",indexHandler)
	http.ListenAndServe(":"+PORT,nil)
}
