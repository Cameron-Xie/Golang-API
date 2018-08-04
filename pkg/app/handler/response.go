package handler

import "net/http"

type Response struct {
	payload   interface{}
	code      int
	isSuccess bool
}

func NewResponse(payload interface{}, code int, isSuccess bool) *Response {
	return &Response{
		payload:   payload,
		code:      code,
		isSuccess: isSuccess,
	}
}

func GetResponse(payload interface{}, err error, code int) *Response {
	if err != nil {
		return ErrorResponse(err)
	}

	return NewResponse(payload, code, true)
}

func ErrorResponse(err error) *Response {
	switch err.(type) {
	case *ResponseError:
		responseErr := err.(*ResponseError)
		return NewResponse(responseErr.Payload(), responseErr.StateCode(), false)
		break
	}

	return NewResponse(err.Error(), http.StatusInternalServerError, false)
}

func writeResponse(w http.ResponseWriter) func([]byte, int) {
	return func(payload []byte, code int) {
		w.WriteHeader(code)
		w.Write(payload)
	}
}
