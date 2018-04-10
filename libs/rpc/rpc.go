package rpc


type DataHandler interface {
	HandleRequest(string) string
}

type Protocol interface {
	Handle()
}

type AckHandler struct {
	prefix string
}

func (hn *AckHandler) Init() {
	if hn.prefix == "" {
		hn.prefix = "lolkek1"
	}
}

func (hn *AckHandler) HandleRequest(req string) string {
	return hn.prefix + req
}


