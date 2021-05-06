package https

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"thingworks.net/thingworks/jarvis-boot/autoconfig/config"
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
)

type HttpHandler func(http.ResponseWriter, *HttpRequest)
type HandlerMap map[string]HttpHandler

func (handlers HandlerMap) GetHandler(key string) (HttpHandler, error) {
	handler := handlers[strings.ToUpper(key)]

	if handler == nil {
		return nil, errors.New(fmt.Sprintf("%s does NOT exist", key))
	}

	return handler, nil
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

const (
	MethodNotAllowed   = "request.method.not.allowed"
	RequestIsInValid   = "request.is.invalid"
	PolicyNotSupported = "policy.not.supported"
	NotAbleToGrant     = "not.able.to.grant"
)

type Response struct {
	Status  int          `json:"status"`
	Message string       `json:"message,omitempty"`
	Code    string       `json:"code,omitempty"`
	Result  interface{}  `json:"result,omitempty"`
	Request *HttpRequest `json:"-"`
}

var ResponseMethodNotAllowed = Response{Status: 400, Code: MethodNotAllowed, Message: "method is NOT allowed"}

func (result Response) To(w http.ResponseWriter) {
	writeToResponse(w, result, true)
}

func (result Response) ToWithoutLog(w http.ResponseWriter) {
	writeToResponse(w, result, false)
}

func writeToResponse(w http.ResponseWriter, result Response, writeLog bool) {
	stringResult := strings2.ToJsonString(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.Status)

	if writeLog {
		auditLog(result, stringResult)
	}

	_, _ = io.WriteString(w, stringResult)
}

func auditLog(result Response, response string) {
	requestLog := ""
	if result.Request != nil {
		requestLog = fmt.Sprintf("\nReceive: [%s] of request [%s]\n", result.Request.Method, result.Request.RequestURI)
	}

	log.Printf("%sResponse: [%s]", requestLog, response)
}

func Get(handler HttpHandler) HttpHandler {
	return parseForm(method(GET, handler))
}

type HttpRequest struct {
	*http.Request
	pathParams map[string]string
}

func newHttpRequest(request *http.Request) *HttpRequest {
	return &HttpRequest{
		Request:    request,
		pathParams: map[string]string{},
	}
}

func (request *HttpRequest) QueryParams() map[string]string {
	result := map[string]string{}

	for key := range request.URL.Query() {
		result[key] = request.GetQueryParam(key)
	}

	return result
}

func (request *HttpRequest) GetQueryParam(query string) string {
	return request.QueryString(query)
}

func (request *HttpRequest) GetPathParam(name string) string {
	return request.pathParams[name]
}

func (request *HttpRequest) Read(buf []byte) (int, error) {
	n, err := request.Body.Read(buf)

	if err != nil {
		return n, err
	}

	return n, request.Body.Close()
}

func (request *HttpRequest) ApiKey() string {
	apiKey := "apiKey"

	for key, values := range request.Header {
		if strings2.EqualCaseIgnored(key, apiKey) && len(values) > 0 {
			return values[0]
		}
	}

	return request.QueryString(apiKey)
}

func (request *HttpRequest) QueryString(key string) string {
	_ = request.ParseForm()

	for queryName, values := range request.Form {
		if strings2.EqualCaseIgnored(queryName, key) {
			if len(values) > 0 {
				return values[0]
			}

			return ""
		}
	}

	return ""
}

func (request *HttpRequest) AddPathParams(pathParams map[string]string) {
	for key, value := range pathParams {
		request.pathParams[key] = value
	}
}

func Post(handler HttpHandler) HttpHandler {
	return parseForm(method(POST, handler))
}

func Mul(handlers HandlerMap) HttpHandler {
	return func(w http.ResponseWriter, r *HttpRequest) {
		handler, err := handlers.GetHandler(r.Method)

		if err != nil {
			ResponseMethodNotAllowed.To(w)
			return
		}

		parseForm(handler)(w, r)
	}
}

func parseForm(handler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *HttpRequest) {
		_ = r.ParseForm()
		handler(w, r)
	}
}

func method(methodName string, handler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *HttpRequest) {
		if !strings2.EqualCaseIgnored(r.Method, methodName) {
			ResponseMethodNotAllowed.To(w)
			return
		}

		handler(w, r)
	}
}

func Register(resource Resource, router *mux.Router, path string) {
	for key, handler := range resource.Handlers() {
		reqPath, _ := getCompleteUrl(path, key)
		router.HandleFunc(reqPath, wrapToFunc(reqPath, handler))
	}
}

func getCompleteUrl(path string, key string) (string, error) {
	if strings2.IsNotBlank(key) {
		return url.PathUnescape(strings2.Join([]string{path, key}, "/"))
	}

	return url.PathUnescape(path)
}

func wrapToFunc(path string, handler HttpHandler) func(w http.ResponseWriter, r *http.Request) {
	AddPermission(path)
	return func(w http.ResponseWriter, r *http.Request) {
		defer processError(w)()
		auth(path, handler)(w, newHttpRequest(r))
	}
}

func processError(w http.ResponseWriter) func() {
	return func() {
		if err := recover(); err != nil {
			logrus.Error(err)

			appErr, ok := err.(ApplicationError)

			if ok {
				err = appErr
			} else {
				err = writeErrorInfo(w, err)
			}

			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func writeErrorInfo(w http.ResponseWriter, err interface{}) interface{} {
	w.WriteHeader(500)
	appError := WrapIntoAppError(err)

	commonErrorStruct := struct {
		Code      int       `json:"code"`
		Message   string    `json:"message"`
		ErrorCode ErrorCode `json:"errorCode"`
	}{
		Code:      appError.Code(),
		Message:   appError.Message(),
		ErrorCode: appError.ErrorCode(),
	}

	data, _ := json.Marshal(commonErrorStruct)

	_, err = w.Write(data)
	return err
}

func auth(path string, handler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *HttpRequest) {

		if permissionSet.Contains(path) {
			if hasPlaceHolder(path) {
				vars := mux.Vars(r.Request)

				r.AddPathParams(vars)
			}

			handler(w, r)
			return
		}

		if r.ApiKey() != config.DefaultConfig().ApiKey {
			Response{Status: 401, Message: "No Permission", Code: "no.permission"}.To(w)
			return
		}

		handler(w, r)
	}
}

func hasPlaceHolder(path string) bool {
	compile, err := regexp.Compile("{.*}")

	if err != nil {
		return false
	}

	return compile.MatchString(path)
}
