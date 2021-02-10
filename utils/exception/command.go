package exception

import "thingworks/common/https"

type IllegalCommandException struct {
	statusCode int
	message string
	errorCode https.ErrorCode
}

func NewIllegalCommandException(statusCode int, message string, errorCode https.ErrorCode) *IllegalCommandException {
	return &IllegalCommandException{statusCode, message, errorCode}
}

func (exception *IllegalCommandException) Code() int {
	return exception.statusCode
}

func (exception *IllegalCommandException) Message() string {
	return exception.message
}

func (exception *IllegalCommandException) ErrorCode() https.ErrorCode {
	return exception.errorCode
}

func (exception *IllegalCommandException) Error() string {
	return exception.Message()
}


