package main

import (
	"fmt"
	"github.com/spf13/viper"
	"sieged/internal/flags"
	"sieged/pkg/helpers"
	"sieged/pkg/random"
	"sieged/pkg/storage"
	"net/http"
	"strconv"
)

var serviceChar = ""

var Storage storage.HsetRadixStorage

func GenerateFlag() string {
	return serviceChar + random.String(30, random.Digits+random.UpperCaseBytes) + "="
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.FormValue("team_id")
	roundId := r.FormValue("round_id")
	fmt.Println(teamId, roundId)
	if teamId == "" || roundId == "" {
		fmt.Fprint(w, "You shuld set round and team id.")
		return
	}
	teamIdInt, _ := strconv.Atoi(teamId)
	roundIdInt, _ := strconv.Atoi(roundId)
	data := flags.Data{Team: teamIdInt, Round: roundIdInt}

	flagJson, _ := flags.DumpData(&data)
	flag := GenerateFlag()
	Storage.Set(flag, flagJson)

	fmt.Fprint(w, flag)
}

func main() {
	viper.SetConfigFile("flag_adder.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("redis_host", "127.0.0.1")
	viper.SetDefault("redis_port", "6379")
	viper.SetDefault("redis_pool_size", 5)
	viper.SetDefault("http_host", "127.0.0.1")
	viper.SetDefault("http_port", "5000")
	err := viper.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err, "Failed to read config")

	serviceChar = viper.GetString("flag_prefix")

	Rp := new(storage.RadixPool)
	Rp.Build(viper.GetString("redis_host"), viper.GetString("redis_port"), viper.GetInt("redis_pool_size"))

	Storage = storage.HsetRadixStorage{RadixP: Rp, SetName: "flags"}
	http.HandleFunc("/", indexHandler)
	httpHost := viper.GetString("http_host")
	httpPort := viper.GetString("http_port")
	fmt.Printf("Flag Adder started on %s:%s \n", httpHost, httpPort)
	err = http.ListenAndServe(httpHost+":"+httpPort, nil)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
