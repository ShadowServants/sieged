package main


import (
	"net/http"
	"fmt"
	"math/rand"
	"time"
	"strconv"
	"github.com/jnovikov/hackforces/libs/storage"
	"github.com/jnovikov/hackforces/libs/flagdata"
)
var (
    PORT = "8011"
    REDIS_POOL_SIZE = 30
    REDIS_PORT = "6379"
    REDIS_IP = "localhost"
)


const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const serviceChar = "R"

var Stor storage.HsetRadixStorage

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
	data := flagdata.FlagData{team_id_int,round_id_int}

	flag_json,_ := flagdata.DumpFlagData(&data)
	flag := GenerateFlag()
	Stor.Set(flag,flag_json)

	fmt.Fprint(w,flag)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	Rp := new(storage.RadixPool)
	Rp.Build("127.0.0.1","6379",7)
	Stor = storage.HsetRadixStorage{Rp,"flags"}
	http.HandleFunc("/",indexHandler)
	http.ListenAndServe("127.0.0.1"+":"+PORT,nil)
}


