package main

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"sieged/internal/flags"
	"sieged/internal/team"
	"sieged/pkg/helpers"
	"sieged/pkg/storage"
	"strconv"
)

type Service struct {
	FPrefix  string
	HostPort string
}

func main() {

	viper.SetConfigFile("router_config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "8000")
	viper.SetDefault("team_num", 10)
	viper.SetDefault("redis_host", "127.0.0.1")
	viper.SetDefault("redis_port", "6379")
	viper.SetDefault("redis_pool_size", 20)
	err := viper.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err, "Failed to read router_config.yaml")

	var services []Service
	err = viper.UnmarshalKey("services", &services)
	helpers.FailOnError(err, "Failed to load services")

	radixPool := new(storage.RadixPool)
	radixPool.Build(viper.GetString("redis_host"),
		viper.GetString("redis_port"), viper.GetInt("redis_pool_size"))

	routerHandler := flags.NewRouter(viper.GetInt("team_num"))

	for _, service := range services {
		print(service.FPrefix)
		routerHandler.RegisterTCPHandler(service.FPrefix, service.HostPort)
	}

	viperTeams := viper.New()
	viperTeams.SetConfigFile("teams.yaml")
	viperTeams.SetConfigType("yaml")
	viperTeams.AddConfigPath(".")

	var teams []team.Address
	err = viperTeams.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err, "Failed to read teams.yaml")

	err = viperTeams.UnmarshalKey("teams", &teams)
	helpers.FailOnError(err, "Cant parse team list")
	for _, t := range teams {
		routerHandler.AddTeam(t.Network, strconv.Itoa(t.Id))
	}

	log.Print("Ip to team -- ", routerHandler.IpStorage)
	host := viper.GetString("host")
	port := viper.GetString("port")
	if viper.IsSet("visualisation_url") {
		routerHandler.SetVisualisation(viper.GetString("visualisation_url"))

	}
	if viper.IsSet("attack_logs") {
		f, err := os.OpenFile(viper.GetString("attack_logs"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Cant opening log file: %v", err.Error())
		}
		defer f.Close()
		routerHandler.SetLogger(f)
	}
	tcpRouter := new(TcpRouter).SetHost(host).SetPort(port).SetRouter(routerHandler)
	tcpRouter.StartPolling()
}
