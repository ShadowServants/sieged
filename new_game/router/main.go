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
	FPrefix string
	HostPort string
}



func main() {

	viper.SetConfigFile("router_config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("http_host","0.0.0.0")
	viper.SetDefault("http_port","8000")
	viper.SetDefault("team_num",10)
	viper.SetDefault("redis_host","127.0.0.1")
	viper.SetDefault("redis_/port","6379")
	viper.SetDefault("redis_pool_size",20)
	err := viper.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err,"Failed to read router_config.yaml")


	var services []Service
    err = viper.UnmarshalKey("services", &services)
    helpers.FailOnError(err,"Failed to load services")

	radix_pool := new(storage.RadixPool)
	radix_pool.Build(
	viper.GetString("redis_host"),
	viper.GetString("redis_port"),
	viper.GetInt("redis_pool_size"))


	router_handler := flag_router.NewFlagRouter(viper.GetInt("team_num"))

	//router_handler.SetIpStorage(&storage.HsetRadixStorage{radix_pool,"player_ip_to_team"})

	for _, service := range services {
		print(service.FPrefix)
		router_handler.RegisterHandler(service.FPrefix,service.HostPort)
	}

	viper_teams := viper.New()
	viper_teams.SetConfigFile("teams.yaml")
	viper_teams.SetConfigType("yaml")
	viper_teams.AddConfigPath(".")

	var teams []team_list.TeamIP
	err = viper_teams.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err,"Failed to read teams.yaml")

	err = viper_teams.UnmarshalKey("teams",&teams)
	helpers.FailOnError(err,"Cant parse team list")
	for _, team := range teams {
		router_handler.AddTeam(team.Network,strconv.Itoa(team.Id))
	}

	fmt.Print("Ip to team -- ",router_handler.IpStorage)
	http_port := viper.GetString("tcp_port")
	http_host := viper.GetString("tcp_host")
	if viper.IsSet("visualisation_url") {
		router_handler.SetVisualisation(viper.GetString("visualisation_url"))

	}
	tcp_router := new(flag_router.TcpRouter).SetHost(http_host).SetPort(http_port).SetRouter(router_handler)
	tcp_router.StartPolling()


}