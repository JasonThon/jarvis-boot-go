package https

import "net/http"

type Gateway struct {
	resourceMapper ResourceMap
	servMux *http.ServeMux
}

func NewGateway() *Gateway {
	return &Gateway{
		resourceMapper: ResourceMap{},
		servMux: http.NewServeMux(),
	}
}

func (router *Gateway) RegisterResource(resourceMap ResourceMap) {
	for root, resource := range resourceMap {
		router.resourceMapper[root] = resource
		Register(resource, router.servMux, root)
	}
}

func (router *Gateway) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.servMux.ServeHTTP(writer, request)
}
