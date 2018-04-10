package main

import (
	"hackforces/libs/storage"
	"hackforces/service_controller/flag_handler"
	"hackforces/libs/rpc"
	"github.com/spf13/viper"
	"hackforces/libs/helpers"
	"fmt"
)

func main() {
	viper.SetConfigFile("flag_handler.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("redis_host","127.0.0.1")
	viper.SetDefault("redis_port","6379")
	viper.SetDefault("redis_pool_size",10)

	viper.SetDefault("tcp_host","127.0.0.1")
	viper.SetDefault("tcp_pool_size",15)

	err := viper.ReadInConfig() // Find and read the config file
	helpers.FailOnError(err,"Failed to read flag_handler.yaml")


	Rp := new(storage.RadixPool)
	Rp.Build(viper.GetString("redis_host"),viper.GetString("redis_port"),viper.GetInt("redis_pool_size"))
	radixFactory := storage.RadixFactory{Pool: Rp}

	flagHandlerFactory := flaghandler.NewFlagHandlerFactory()
	flagHandlerFactory.SetPointStorage(radixFactory.GetHsetStorage("points"))
	flagHandlerFactory.SetFlagStorage(radixFactory.GetHsetStorage("flags"))
	flagHandlerFactory.SetTeamFlagsSet(radixFactory.GetKeySet())
	flagHandlerFactory.SetRoundStorage(radixFactory.GetHsetStorage("rounds"))
	flagHandlerFactory.SetStatusStorage(radixFactory.GetHsetStorage("statuses"))
	flagHandlerFactory.SetTeamNum(viper.GetInt("team_num"))
	flagHandler := flagHandlerFactory.GetFlagHandler()
	tcpPort := viper.GetString("tcp_port")
	tcpHost := viper.GetString("tcp_host")
	rpcTcp := rpc.TcpRpc{Port: tcpPort, Addr: tcpHost,
							Handler: flagHandler}
	fmt.Printf("Flag Handler started on %s:%s \n", tcpHost, tcpPort)

	rpcTcp.Handle()

}
