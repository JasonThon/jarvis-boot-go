package https

import "encoding/json"

type ErrorCode string

const (
	InternalError ErrorCode = "INTERNAL_ERROR"
	InvalidProperty ErrorCode = "INVALID_PROPERTY"
)

type ApplicationError interface {
	Code() int
	Message() string
	ErrorCode() ErrorCode
}

type internalError struct {
	msg string
}

func (internal *internalError) Code() int {
	return 500
}

func (internal *internalError) Message() string {
	return internal.msg
}

func (internal *internalError) ErrorCode() ErrorCode {
	return InternalError
}


func WrapIntoAppError(err interface{}) ApplicationError {
	bytes, _ := json.Marshal(err)
	return &internalError{
		msg: string(bytes),
	}
}