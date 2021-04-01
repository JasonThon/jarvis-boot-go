package https

import (
	"net/http"
	"thingworks.net/thingworks/common/https/nio"
)

type WrappedRequest struct {
	writer  *nio.BufferedResponseWriter
	request *http.Request
	reqId   uint64
}

func (w *WrappedRequest) ReqId() uint64 {
	return w.reqId
}

func (w *WrappedRequest) SetReqId(reqId uint64) {
	w.reqId = reqId
}

func (w *WrappedRequest) Request() *http.Request {
	return w.request
}

func (w *WrappedRequest) Writer() *nio.BufferedResponseWriter {
	return w.writer
}

func NewWrappedRequest(writer http.ResponseWriter, newReq *http.Request, router GatewayAdaptor) *WrappedRequest {
	return &WrappedRequest{
		writer:  nio.NewBufferedResponseWriter(writer),
		request: newReq,
		reqId:   *router.ReqId(),
	}
}
