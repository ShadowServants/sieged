package main

import (
	"github.com/jnovikov/hackforces/back/service_controller/flag_handler"
	"github.com/jnovikov/hackforces/back/libs/storage"
	"github.com/jnovikov/hackforces/back/libs/flagstorage"
	"github.com/jnovikov/hackforces/back/libs/statusstorage"
	"github.com/jnovikov/hackforces/back/libs/rpc"
)

func main() {
	fl := flaghandler.FlagHandler{}
	Rp := new(storage.RadixPool)
	Rp.Build("localhost","6379",10)
	//executor := storage.GetRedisExecutor("6379",10)
	//defer executor.Close()
	ps := storage.HsetRadixStorage{Rp,"points"}

	//ps := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"points"}

	pointStorage := flaghandler.PointsStorage{&ps}
	fs := storage.HsetRadixStorage{Rp,"flags"}
	//fs := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"flags"}
	flags := flagstorage.NewFlagStorage(&fs)
	ks := storage.RadixKeySet{Rp}
	fl.TeamFlagsSet = &ks
	fl.Points = &pointStorage
	rs := storage.HsetRadixStorage{Rp,"rounds"}
	//rs := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"rounds"}
	roundst := &flaghandler.RoundStorage{&rs}
	fl.RoundSt = roundst
	fl.Flags = flags
	fl.RoundDelta = 3
	fl.RoundCached = false
	ss := storage.HsetRadixStorage{Rp,"statuses"}
	//ss := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"statuses"}
	statuses := statusstorage.NewStatusStorage(&ss)
	fl.StatusStorage = statuses
	fl.Build()
	rpc_tcp := rpc.TcpRpc{"7878","localhost",&fl}
	rpc_tcp.Handle()

}
