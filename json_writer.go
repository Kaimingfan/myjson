package myjson

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	leftParenthesis  = '{'
	rightParenthesis = '}'
	comma            = ','
	quota            = '"'
)

type jsonWriter struct {
	builder *strings.Builder
}

func newJsonWriter() *jsonWriter {
	jw := &jsonWriter{
		builder: &strings.Builder{},
	}
	return jw
}
func (jw *jsonWriter) start() {
	jw.builder.WriteRune(leftParenthesis)
}
func (jw *jsonWriter) end() {
	jw.builder.WriteRune(rightParenthesis)
}
func (jw *jsonWriter) writeKey(key string) {
	jw.builder.WriteString(fmt.Sprintf("\"%s\":", key))
}
func (jw *jsonWriter) writeVal(val interface{}) {
	jw.builder.WriteString(fmt.Sprintf("%v", val))
}
func (jw *jsonWriter) writeString(val interface{}) {
	// 处理换行等特殊符号
	val = strings.Replace(val.(string), "\"", "\\\"", -1)
	val = strings.Replace(val.(string), "\r", "\\r", -1)
	val = strings.Replace(val.(string), "\n", "\\n", -1)
	jw.builder.WriteString(fmt.Sprintf("\"%s\"", val))
}
func (jw *jsonWriter) writeNil() {
	jw.builder.WriteString("null")
}
func (jw *jsonWriter) writeUnknown() {
	jw.builder.WriteString("unknown")
}
func (jw *jsonWriter) lineBreak() {
	jw.builder.WriteRune('\n')
}
func (jw *jsonWriter) comma() {
	jw.builder.WriteRune(',')
}
func (jw *jsonWriter) leftBracket() {
	jw.builder.WriteRune('[')
}
func (jw *jsonWriter) rightBracket() {
	jw.builder.WriteRune(']')
}

func (jw *jsonWriter) WriteObject(rv reflect.Value) {
	if rv.Interface() == nil {
		jw.writeNil()
		return
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		jw.writeString(rv.Interface())
	case reflect.Int, reflect.Int16, reflect.Int64, reflect.Int32, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float64, reflect.Float32, reflect.Bool:
		jw.writeVal(rv.Interface())
	case reflect.Slice:
		jw.leftBracket()
		for i := 0; i < rv.Len(); i++ {
			if i != 0 {
				jw.comma()
			}
			jw.WriteObject(rv.Index(i))
		}
		jw.rightBracket()
	case reflect.Struct:
		jw.start()
		// do something
		for i := 0; i < rv.NumField(); i++ {
			key := rv.Type().Field(i).Tag.Get("json")
			if key == "" {
				key = rv.Type().Field(i).Name
			}
			subVal := rv.Field(i)
			jw.writeKey(key)
			jw.WriteObject(subVal)
			if i < rv.NumField()-1 {
				jw.comma()
			}
			// jw.lineBreak()
		}
		jw.end()
	case reflect.Interface:
		rv = rv.Elem()
		jw.WriteObject(rv)
	default:
		jw.writeUnknown()
	}
}

func (jw *jsonWriter) String() string {
	return jw.builder.String()
}
