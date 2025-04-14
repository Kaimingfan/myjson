package myjson

import (
	"reflect"

	"github.com/pkg/errors"
)

func MarshalString(v interface{}) (string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return "", errors.New("invalid input")
	}
	jw := newJsonWriter()
	jw.WriteObject(rv)
	return jw.String(), nil
}
