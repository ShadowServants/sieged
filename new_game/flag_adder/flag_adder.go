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
	teamId := r.FormValue("team_id")
	roundId := r.FormValue("round_id")
	fmt.Println(teamId, roundId)
	if teamId == "" || roundId == ""{
		fmt.Fprint(w,"You shuld set round and team id.")
		return
	}
	teamIdInt,_ := strconv.Atoi(teamId)
	roundIdInt, _ := strconv.Atoi(roundId)
	data := flagdata.FlagData{Team: teamIdInt, Round: roundIdInt}

	flagJson,_ := flagdata.DumpFlagData(&data)
	flag := GenerateFlag()
	Stor.Set(flag, flagJson)

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

	Stor = storage.HsetRadixStorage{RadixP: Rp, SetName: "flags"}
	http.HandleFunc("/",indexHandler)
	httpHost := viper.GetString("http_host")
	httpPort := viper.GetString("http_port")
	fmt.Printf("Flag Adder started on %s:%s \n", httpHost, httpPort)
	err = http.ListenAndServe(httpHost+":"+httpPort,nil)
	if err != nil {
		fmt.Printf(err.Error())
	}
}




