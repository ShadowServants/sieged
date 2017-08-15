package main

import (
	"net/http"
	"fmt"
	"github.com/jnovikov/hackforces/back/libs/storage"
	"github.com/jnovikov/hackforces/back/flag_router"
)

type HTTPFlagRouter struct {
	Fr *flag_router.FlagRouter
	Port string
}

func (fr *HTTPFlagRouter) SetPort(port string) *HTTPFlagRouter {
	fr.Port = port
	return fr
}

func (fr *HTTPFlagRouter) SetRouter(router *flag_router.FlagRouter) *HTTPFlagRouter {
	fr.Fr = router
	return fr
}

func (fr *HTTPFlagRouter) handleRequest(w http.ResponseWriter, r *http.Request) {
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

func (fh *HTTPFlagRouter) StartPolling() {
	http.HandleFunc("/",fh.handleRequest)
	http.ListenAndServe("0.0.0.0:"+fh.Port,nil)
}

func main() {
	Rp := new(storage.RadixPool)
	Rp.Build("127.0.0.1","6379",20)
	FlagRouter := flag_router.NewFlagRouter(7)
	FlagRouter.SetVisualisation("http://127.0.0.1:3000/broadcast")
	FlagRouter.IpStorage = &storage.HsetRadixStorage{Rp,"player_ip_to_team"}
	FlagRouter.RegisterHandler("T","127.0.0.1:8012")
	FlagRouter.RegisterHandler("C","127.0.0.1:7012")
	HttpRouter := new(HTTPFlagRouter).SetPort("7331").SetRouter(FlagRouter)
	HttpRouter.StartPolling()
}
