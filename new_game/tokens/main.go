package main

import (
	"github.com/spf13/viper"
	"hackforces/libs/team_list"
	"hackforces/libs/helpers"
	"flag"
	"hackforces/libs/storage"
	"hackforces/libs/token"
	"fmt"
)

func main() {
	viper.SetConfigFile("router_config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("redis_host", "127.0.0.1")
	viper.SetDefault("redis_port", "6379")
	viper.SetDefault("redis_pool_size", 20)
	viper.ReadInConfig()

	generate := flag.String("generate", "", "That command ill regenerate new tokens for teams. WARNING! All tokens will be lost")

	flag.Parse()

	viperTeams := viper.New()
	viperTeams.SetConfigFile("teams.yaml")
	viperTeams.SetConfigType("yaml")
	viperTeams.AddConfigPath(".")

	Rp := new(storage.RadixPool)
	Rp.Build(viper.GetString("redis_host"), viper.GetString("redis_port"), viper.GetInt("redis_pool_size"))
	radixFactory := storage.RadixFactory{Pool: Rp}

	ts := token.NewStorage(radixFactory.GetHsetStorage("tokens"))

	var teams []team_list.TeamIP
	err := viperTeams.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err, "Failed to read teams.yaml")

	err = viperTeams.UnmarshalKey("teams", &teams)
	helpers.FailOnError(err, "Cant parse team list")
	if *generate != "" {
		for _, team := range teams {
			ts.New(team.Id)
		}
	}

	for _, team := range teams {
		s, err := ts.FindById(team.Id)
		helpers.FailOnError(err, "Failed to get token by id")
		fmt.Printf("%d. %s: %s \n", team.Id, team.Name, s.Token)
	}
}
