package main

import (
	"net/http"
	"fmt"
	"math/rand"
	"time"
	"strconv"
	"hackforces/libs/storage"
	"hackforces/libs/flagdata"
	"github.com/spf13/viper"
	"hackforces/libs/helpers"
)


const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var serviceChar = ""

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
	viper.SetConfigFile("flag_adder.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("redis_host","127.0.0.1")
	viper.SetDefault("redis_port","6379")
	viper.SetDefault("redis_pool_size",5)
	viper.SetDefault("http_host", "127.0.0.1")
	viper.SetDefault("http_port", "5000")
	err := viper.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err,"Failed to read config")

	serviceChar = viper.GetString("flag_prefix")

	rand.Seed(time.Now().UnixNano())
	Rp := new(storage.RadixPool)
	Rp.Build(viper.GetString("redis_host"),viper.GetString("redis_port"),viper.GetInt("redis_pool_size"))

	Stor = storage.HsetRadixStorage{Rp,"flags"}
	http.HandleFunc("/",indexHandler)
	http_host := viper.GetString("http_host")
	http_port := viper.GetString("http_port")
	fmt.Printf("Flag Adder started on %s:%s \n", http_host, http_port)
	http.ListenAndServe(http_host+":"+http_port,nil)
}




