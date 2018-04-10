package main

import (
	"hackforces/libs/storage"
	"hackforces/service_controller/round_handler"
	"hackforces/libs/rpc"
	"github.com/spf13/viper"
	"hackforces/libs/helpers"
	"hackforces/libs/team_list"
)


func main(){

	viper.SetConfigFile("round_handler.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("redis_host","127.0.0.1")
	viper.SetDefault("redis_port","6379")
	viper.SetDefault("redis_pool_size",10)

	viper.SetDefault("http_host","127.0.0.1")
	viper.SetDefault("http_port","9090")
	viper.SetDefault("default_points",1700)
	err := viper.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err,"Failed to read flag_handler.yaml")

	Rp := new(storage.RadixPool)
	Rp.Build(viper.GetString("redis_host"),viper.GetString("redis_port"),viper.GetInt("redis_pool_size"))
	radixFactory := storage.RadixFactory{Pool: Rp}

	factory := round_handler.NewHandlerFactory()
	factory.SetIpStorage(radixFactory.GetHsetStorage("team_to_ip"))
	factory.SetPointStorage(radixFactory.GetHsetStorage("points"))
	factory.SetRoundStorage(radixFactory.GetHsetStorage("rounds"))
	factory.SetStatusStorage(radixFactory.GetHsetStorage("statuses"))
	factory.SetTeamStorage(radixFactory.GetHsetStorage("teams_id"))

	rHandler := factory.GetHandler()
	viperTeams := viper.New()
	viperTeams.SetConfigFile("teams.yaml")
	viperTeams.SetConfigType("yaml")
	viperTeams.AddConfigPath(".")

	var teams []team_list.TeamIP
	err = viperTeams.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err,"Failed to read teams.yaml")

	err = viperTeams.UnmarshalKey("teams",&teams)
	helpers.FailOnError(err,"Cant parse team list")

	rHandler.LoadsTeamsIp(teams)
	rHandler.DefaultPoints = viper.GetInt("default_points")
	rHandler.CheckerName = viper.GetString("checker_name")
	server := rpc.NewRpcServer(viper.GetString("http_host"),viper.GetString("http_port"))

	server.Register("/round", rHandler)
	server.Handle()







}
