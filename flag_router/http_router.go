package flag_router

import (
	"net/http"
	"fmt"
)

type HTTPFlagRouter struct {
	Fr *FlagRouter
	Port string
	Host string
}

func (fr *HTTPFlagRouter) SetPort(port string) *HTTPFlagRouter {
	fr.Port = port
	return fr
}

func (fr *HTTPFlagRouter) SetHost(host string) *HTTPFlagRouter {
	fr.Host = host
	return fr
}

func (fr *HTTPFlagRouter) SetRouter(router *FlagRouter) *HTTPFlagRouter {
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
	http.ListenAndServe(fh.Host+":"+fh.Port,nil)
}