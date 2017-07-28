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
	Rp.Build("127.0.0.1","6378",25)
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
	fl.RoundDelta = 5
	fl.RoundCached = false
	ss := storage.HsetRadixStorage{Rp,"statuses"}
	//ss := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"statuses"}
	statuses := statusstorage.NewStatusStorage(&ss)
	fl.StatusStorage = statuses
	fl.Build()
	rpc_tcp := rpc.TcpRpc{"8012","127.0.0.1",&fl}
	rpc_tcp.Handle()

}
