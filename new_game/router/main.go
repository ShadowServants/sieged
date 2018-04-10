package main

import (
	"hackforces/flag_router"
	"github.com/spf13/viper"
	"hackforces/libs/storage"
	"hackforces/libs/helpers"
	"hackforces/libs/team_list"
	"strconv"
	"fmt"
)

type Service struct {
	FPrefix  string
	HostPort string
}

func main() {

	viper.SetConfigFile("router_config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("http_host", "0.0.0.0")
	viper.SetDefault("http_port", "8000")
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

	routerHandler := flag_router.NewFlagRouter(viper.GetInt("team_num"))

	//router_handler.SetIpStorage(&storage.HsetRadixStorage{radix_pool,"player_ip_to_team"})

	for _, service := range services {
		print(service.FPrefix)
		routerHandler.RegisterHandler(service.FPrefix, service.HostPort)
	}

	viperTeams := viper.New()
	viperTeams.SetConfigFile("teams.yaml")
	viperTeams.SetConfigType("yaml")
	viperTeams.AddConfigPath(".")

	var teams []team_list.TeamIP
	err = viperTeams.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err, "Failed to read teams.yaml")

	err = viperTeams.UnmarshalKey("teams", &teams)
	helpers.FailOnError(err, "Cant parse team list")
	for _, team := range teams {
		routerHandler.AddTeam(team.Network, strconv.Itoa(team.Id))
	}

	fmt.Print("Ip to team -- ", routerHandler.IpStorage)
	httpPort := viper.GetString("tcp_port")
	httpHost := viper.GetString("tcp_host")
	if viper.IsSet("visualisation_url") {
		routerHandler.SetVisualisation(viper.GetString("visualisation_url"))

	}
	tcpRouter := new(flag_router.TcpRouter).SetHost(httpHost).SetPort(httpPort).SetRouter(routerHandler)
	tcpRouter.StartPolling()

}
