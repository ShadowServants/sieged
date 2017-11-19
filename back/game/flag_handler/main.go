package main

import (
	"github.com/jnovikov/hackforces/back/service_controller/flag_handler"
	"github.com/jnovikov/hackforces/back/libs/storage"
	//"github.com/jnovikov/hackforces/back/libs/flagstorage"
	//"github.com/jnovikov/hackforces/back/libs/statusstorage"
	"github.com/jnovikov/hackforces/back/libs/rpc"
)

func main() {
	Rp := new(storage.RadixPool)
	Rp.Build("127.0.0.1","6379",25)
	radixFactory := storage.RadixFactory{Rp}

	flagHandlerFactory := flaghandler.NewFlagHandlerFactory()
	flagHandlerFactory.SetPointStorage(radixFactory.GetHsetStorage("points"))
	flagHandlerFactory.SetFlagStorage(radixFactory.GetHsetStorage("flags"))
	flagHandlerFactory.SetTeamFlagsSet(radixFactory.GetKeySet())
	flagHandlerFactory.SetRoundStorage(radixFactory.GetHsetStorage("rounds"))
	flagHandlerFactory.SetStatusStorage(radixFactory.GetHsetStorage("statuses"))
	flagHandler := flagHandlerFactory.GetFlagHandler()
	rpc_tcp := rpc.TcpRpc{"7878", "127.0.0.1", flagHandler}
	rpc_tcp.Handle()

	//fl := flaghandler.FlagHandler{}
	//executor := storage.GetRedisExecutor("6379",10)
	//defer executor.Close()
	//ps := storage.HsetRadixStorage{Rp,"points"}

	//ps := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"points"}

	//pointStorage := flaghandler.PointsStorage{&ps}
	//fs := storage.HsetRadixStorage{Rp,"flags"}
	//fs := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"flags"}
	//flags := flagstorage.NewFlagStorage(&fs)
	//ks := storage.RadixKeySet{Rp}
	//fl.TeamFlagsSet = &ks
	//fl.Points = &pointStorage
	//rs := storage.HsetRadixStorage{Rp,"rounds"}
	//rs := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"rounds"}
	//roundst := &flaghandler.RoundStorage{&rs}
	//fl.RoundSt = roundst
	//fl.Flags = flags
	//ss := storage.HsetRadixStorage{Rp,"statuses"}
	//ss := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"statuses"}
	//statuses := statusstorage.NewStatusStorage(&ss)
	//fl.StatusStorage = statuses
	//fl.SetSt
	//fl.Build()

}
