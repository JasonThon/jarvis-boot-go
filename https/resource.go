package https

type ResourceMap map[string]Resource

type Resource interface {
	Handlers() HandlerMap
}
