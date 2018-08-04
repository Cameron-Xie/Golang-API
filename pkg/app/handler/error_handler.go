package handler

type ErrorHandler func(err []error) *ResponseError

func JsonErrorHandler(stateCode int) ErrorHandler {
	return func(err []error) *ResponseError {
		return NewJsonResponseError(stateCode, err)
	}
}
