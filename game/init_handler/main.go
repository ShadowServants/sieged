package main

import (
	"github.com/jnovikov/hackforces/service_controller/init_handler"
	"github.com/jnovikov/hackforces/service_controller/flag_handler"
	"github.com/jnovikov/hackforces/libs/storage"
	"github.com/jnovikov/hackforces/libs/rpc"
	"github.com/jnovikov/hackforces/service_controller/round_handler"
	"sync"
	"github.com/jnovikov/hackforces/libs/statusstorage"
)

//BAse redis
//func main(){
//	executor := storage.GetRedisExecutor("6379",30)
//	defer executor.Close()
//
//	ps := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"points"}
//	point_st := flaghandler.PointsStorage{&ps}
//	ts := storage.SimpleRedisStorage{storage.BaseRedisStorage{executor}}
//	ih := init_handler.InitHandler{&point_st,&ts}
//
//	rh := new(round_handler.RoundHandler)
//	rh.TeamStorage = &ts
//	rh.CheckerName = "simple.py"
//	rh.IpStorage = &storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"servers_ip"}
//	rs := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"rounds"}
//	rh.Rounds = &flaghandler.RoundStorage{&rs}
//	rh.Wg = sync.WaitGroup{}
//	rh.TeamIds = make([]int,0)
//	ss := storage.HsetRedisStorage{storage.BaseRedisStorage{executor},"statuses"}
//	rh.St = statusstorage.NewStatusStorage(&ss)
//	server := rpc.NewRpcServer("localhost","8081")
//	server.Register("/init",&ih)
//	server.Register("/round",rh)
//	server.Handle()
//}

func main(){

	Rp := new(storage.RadixPool)
	Rp.Build("127.0.0.1","6379",15)
	ps := storage.HsetRadixStorage{Rp,"points"}
	point_st := flaghandler.PointsStorage{&ps}
	ts := storage.SimpleRadixStorage{Rp}
	ih := init_handler.InitHandler{&point_st,&ts}

	rh := new(round_handler.RoundHandler)
	rh.TeamStorage = &ts
	rh.Points = &point_st
	rh.CheckerName = "roskom_check.py"
	rh.IpStorage = &storage.HsetRadixStorage{Rp,"team_to_ip"}
	rs := storage.HsetRadixStorage{Rp,"rounds"}
	rh.Rounds = &flaghandler.RoundStorage{&rs}
	rh.Wg = sync.WaitGroup{}
	rh.TeamIds = make([]int,0)
	ss := storage.HsetRadixStorage{Rp,"statuses"}
	rh.St = statusstorage.NewStatusStorage(&ss)
	server := rpc.NewRpcServer("127.0.0.1","8009")
	server.Register("/init",&ih)
	server.Register("/round",rh)
	server.Handle()
}
