package myjson

import (
	v1 "myjson/v1"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

var (
	bytesPool = sync.Pool{}
)

func newBytes() []byte {
	if ret := bytesPool.Get(); ret != nil {
		return ret.([]byte)
	} else {
		return make([]byte, 0, 128*1024)
	}
}
func freeBytes(p []byte) {
	p = p[:0]
	bytesPool.Put(p)
}

func Marshal(v interface{}) ([]byte, error) {
	buf := newBytes()
	// encode into buf
	err := encodeInto(v, &buf)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, len(buf))
	copy(ret, buf)
	/* return the buffer into pool */
	freeBytes(buf)
	return ret, nil
}

func encodeInto(v interface{}, buf *[]byte) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errors.New("invalid input")
	}
	*buf = v1.WriteObject(*buf, rv)
	return nil
}
