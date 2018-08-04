package functional

import (
	"errors"
	"reflect"
)

type ComposeFunc func(interface{}) (interface{}, error)

func validateFunc(fnI interface{}) error {
	fn := reflect.ValueOf(fnI)

	if fn.Kind() != reflect.Func {
		return errors.New("require a function")
	}

	fnT := fn.Type()

	if fnT.NumIn() != 1 ||
		fnT.NumOut() != 2 ||
		fnT.Out(1).String() != "error" {
		return errors.New("require a function with one input and return two output, second output is error type")
	}

	return nil
}

func Compose(fns ...interface{}) func(interface{}) (interface{}, error) {
	return func(i interface{}) (interface{}, error) {
		var res []reflect.Value

		for _, fnI := range fns {
			if err := validateFunc(fnI); err != nil {
				return nil, err
			}

			fn := reflect.ValueOf(fnI)

			if res == nil {
				res = fn.Call([]reflect.Value{reflect.ValueOf(i)})
			} else {
				res = fn.Call([]reflect.Value{res[0]})
			}

			if res[1].Interface() != nil {
				return nil, res[1].Interface().(error)
			}
		}

		return res[0].Interface(), nil
	}
}
