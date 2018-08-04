package handler

import (
	"strings"
)

type ResponseError struct {
	payload   []error
	stateCode int
}

func NewJsonResponseError(stateCode int, payload []error) *ResponseError {
	return &ResponseError{
		payload:   payload,
		stateCode: stateCode,
	}
}

func (e *ResponseError) Error() string {
	var str strings.Builder

	for _, err := range e.payload {
		str.WriteString(err.Error() + " ")
	}

	return str.String()
}

func (e *ResponseError) Payload() []string {
	var strAry []string

	for _, err := range e.payload {
		strAry = append(strAry, err.Error())
	}

	return strAry
}

func (e *ResponseError) StateCode() int {
	return e.stateCode
}
