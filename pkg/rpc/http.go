package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpRpcServer struct {
	Port       string
	Host       string
	mux        *http.ServeMux
	handlerMap map[string]DataHandler
}

func NewRpcServer(host string, port string  ) *HttpRpcServer {
	ht := new(HttpRpcServer)
	ht.Port = port
	ht.Host = host
	ht.Build()
	return ht
}

type Request struct {
	Request string
}



func (rs *HttpRpcServer) Register(url string,handler DataHandler) {
	rs.handlerMap[url] = handler
}

func (rs *HttpRpcServer) HandleHTTP(w http.ResponseWriter,r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	uri := r.RequestURI
	request := req.Request
	hndlr, _ := rs.handlerMap[uri]
	if err != nil {
		http.Error(w,err.Error(),404)
		return
	}
	if hndlr == nil {
	 	http.Error(w,err.Error(),404)
		return
	}
	fmt.Fprint(w,hndlr.HandleRequest(request))
}


func (rs *HttpRpcServer) Build() {
	rs.handlerMap = make(map[string]DataHandler)
	rs.mux = http.NewServeMux()
	rs.mux.HandleFunc("/",rs.HandleHTTP)
}

func (rs *HttpRpcServer) Handle() {
	http.ListenAndServe(rs.Host+":"+rs.Port, rs.mux)
}






