package main

import (
	"fmt"
	"net/http"
	"sieged/internal/flags"
	"sieged/internal/team/token"
	"sieged/pkg/storage"
)

type HTTPFlagRouter struct {
	Fr   *flags.Router
	ts   *token.Storage
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

func (fr *HTTPFlagRouter) SetRouter(router *flags.Router) *HTTPFlagRouter {
	fr.Fr = router
	return fr
}

func (fr *HTTPFlagRouter) SetTokenStorage(st storage.Storage) *HTTPFlagRouter {
	fr.ts = token.NewStorage(st)
	return fr
}

func (fr *HTTPFlagRouter) handleRequest(w http.ResponseWriter, r *http.Request) {
	flag := r.FormValue("flag")
	tokenString := r.FormValue("token")
	if flag == "" || tokenString == "" {
		fmt.Fprint(w, "Fields 'flag' and 'token' are required")
		return
	}

	tok, err := fr.ts.Find(tokenString)
	if err != nil {
		fmt.Fprintf(w, "Cant find your team by token with error %s", err.Error())
		return
	}

	response := fr.Fr.HandleFlag(flag, tok.TeamId)
	fmt.Fprint(w, response)
	return
}

func (fr *HTTPFlagRouter) StartPolling() {
	http.HandleFunc("/", fr.handleRequest)
	http.ListenAndServe(fr.Host+":"+fr.Port, nil)
}
