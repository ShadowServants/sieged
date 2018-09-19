package main

import (
	"hackforces/flag_router"
	"github.com/spf13/viper"
	"hackforces/libs/storage"
	"hackforces/libs/helpers"
	"hackforces/libs/team_list"
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

	Rp := new(storage.RadixPool)
	Rp.Build(viper.GetString("redis_host"), viper.GetString("redis_port"), viper.GetInt("redis_pool_size"))
	radixFactory := storage.RadixFactory{Pool: Rp}

	routerHandler := flag_router.NewFlagRouter(viper.GetInt("team_num"))

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

	httpPort := viper.GetString("http_port")
	httpHost := viper.GetString("http_host")

	if viper.IsSet("visualisation_url") {
		routerHandler.SetVisualisation(viper.GetString("visualisation_url"))
	}

	httpRouter := new(flag_router.HTTPFlagRouter).SetHost(httpHost).SetPort(httpPort).SetRouter(routerHandler).SetTokenStorage(radixFactory.GetHsetStorage("tokens"))
	httpRouter.StartPolling()

}
