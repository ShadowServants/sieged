package rpc

import (
	"net/http"
	"encoding/json"
	"fmt"
)

//type HelloArgs struct {
//    Who string
//}
//
//type HelloReply struct {
//    Message string
//}
//
//type HelloService struct {}
//
//func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
//    reply.Message = "Hello, " + args.Who + "!"
//    return nil
//}

type HttpRpcServer struct {
	Port string
	Host string
	mux *http.ServeMux
	handler_map map[string]DataHandler
}

func NewRpcServer(host string, port string  ) *HttpRpcServer {
	ht := new(HttpRpcServer)
	ht.Port = port
	ht.Host = host
	ht.Build()
	return ht
}

type RpcRequest struct {
	Request string
}



func (rs *HttpRpcServer) Register(url string,handler DataHandler) {
	rs.handler_map[url] = handler
}

func (rs *HttpRpcServer) HandleHTTP(w http.ResponseWriter,r *http.Request) {
	var req RpcRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	uri := r.RequestURI
	request := req.Request
	hndlr, _ := rs.handler_map[uri]
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
	rs.handler_map = make(map[string]DataHandler)
	rs.mux = http.NewServeMux()
	rs.mux.HandleFunc("/",rs.HandleHTTP)
}

func (rs *HttpRpcServer) Handle() {
	http.ListenAndServe(rs.Host+":"+rs.Port, rs.mux)
}






