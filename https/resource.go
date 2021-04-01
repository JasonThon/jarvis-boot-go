package https

import "net/http"

type ResourceMap map[string]Resource

type Resource interface {
	Handlers() HandlerMap
}

type GatewayAdaptor interface {
	http.Handler
	RequestQueue() chan *WrappedRequest
	Done() chan int
	ReqId() *uint64
	RegisterResource(resourceMap ResourceMap)
	Close()
	Start()
}
