package main

import (
	"net/http"
	"fmt"
	"github.com/jnovikov/hackforces/back/libs/storage"
	"github.com/jnovikov/hackforces/back/flag_router"
)

type FlagRouter struct {
	Fr *flag_router.FlagRouter
	Port string
}


func (fr *FlagRouter) handleRequest(w http.ResponseWriter, r *http.Request) {
	flag := r.FormValue("flag")
	if flag == "" {
		fmt.Fprint(w,"Field flag is required")
		return
	}
	ip := r.RemoteAddr
	response := fr.Fr.HandleRequest(flag,ip)
	fmt.Println(ip,flag,response)
	fmt.Fprint(w,response)
	return
}

func (fh *FlagRouter) StartPolling() {
	http.HandleFunc("/",fh.handleRequest)
	http.ListenAndServe("0.0.0.0:"+fh.Port,nil)
}

func main() {
	Rp := new(storage.RadixPool)
	Rp.Build("127.0.0.1","6379",20)
	fr := new(FlagRouter)
	fr.Fr = flag_router.NewFlagRouter(7)
	fr.Fr.VisualisationEnabled = true
	fr.Fr.VisualisationUrl = "http://127.0.0.1:3000/broadcast"
	fr.Fr.IpStorage = &storage.HsetRadixStorage{Rp,"player_ip_to_team"}
	fr.Fr.RegisterHandler("T","127.0.0.1:8012")
	fr.Fr.RegisterHandler("C","127.0.0.1:7012")
	fr.Port = "7331"
	fr.StartPolling()
}
