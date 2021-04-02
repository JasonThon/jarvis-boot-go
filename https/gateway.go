package https

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"sync/atomic"
	"thingworks.net/thingworks/common/https/nio"
)

const maxRequestNum = 2147483648

var wg sync.WaitGroup

var reqSignal sync.Map

type Gateway struct {
	servMux  *mux.Router
	queue    chan *WrappedRequest
	complete chan int
	reqId    uint64
}

func (router *Gateway) RequestQueue() chan *WrappedRequest {
	return router.queue
}

func (router *Gateway) Done() chan int {
	return router.complete
}

func (router *Gateway) ReqId() *uint64 {
	return &router.reqId
}

func NewGateway() *Gateway {
	g := &Gateway{
		servMux:  mux.NewRouter(),
		queue:    make(chan *WrappedRequest, maxRequestNum),
		reqId:    0,
		complete: make(chan int),
	}

	return g
}

func (router *Gateway) RegisterResource(resourceMap ResourceMap) {
	for root, resource := range resourceMap {
		Register(resource, router.servMux, root)
	}
}

func (router *Gateway) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	newReq, err := nio.NewBufferedRequest(request)

	if err != nil {
		logrus.Error(err)
		return
	}

	wrapped := NewWrappedRequest(writer, newReq, router)

	AddToReqQueue(router, wrapped)
}

func (router *Gateway) Start() {
	go wg.Wait()

	for {
		select {
		case wrapped := <-router.RequestQueue():
			wg.Add(1)

			go func() {
				router.servMux.ServeHTTP(wrapped.Writer(), wrapped.Request())

				copyAndSendComplete(wrapped)

				wg.Done()
			}()

		case <-router.Done():
			break
		}

	}
}

func (router *Gateway) Close() {
	logrus.Info("Gateway start closing...")
	router.Done() <- 1

	wg.Wait()

	logrus.Info("Gateway has closed")
}

func AddToReqQueue(router GatewayAdaptor, wrapped *WrappedRequest) {
	signal := push(router, wrapped)

	wait(signal)
}

func wait(signal chan int) {
	for {
		select {
		case <-signal:
			close(signal)
			return
		default:
			continue
		}
	}
}

func push(router GatewayAdaptor, wrapped *WrappedRequest) chan int {
	signal := make(chan int)
	updateReqId(router)

	newReqId := *router.ReqId()

	reqSignal.Store(newReqId, signal)
	wrapped.SetReqId(newReqId)

	router.RequestQueue() <- wrapped

	return signal
}

func updateReqId(router GatewayAdaptor) {
	atomic.CompareAndSwapUint64(router.ReqId(), *router.ReqId(), atomic.AddUint64(router.ReqId(), 1))
}

func copyAndSendComplete(wrapped *WrappedRequest) {
	_, err := wrapped.Writer().Copy()

	if err != nil {
		logrus.Errorf("Error when copy data from buffer %v", err)
	}

	value, ok := reqSignal.LoadAndDelete(wrapped.ReqId())

	if ok {
		value.(chan int) <- 1
	}
}
