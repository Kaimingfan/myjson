package v1

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	leftParenthesis  = '{'
	rightParenthesis = '}'
)

func start(input []byte) []byte {
	input = append(input, leftParenthesis)
	return input
}
func end(input []byte) []byte {
	input = append(input, rightParenthesis)
	return input
}
func writeKey(input []byte, key string) []byte {
	input = append(input, '"')
	input = append(input, key...)
	input = append(input, '"')
	input = append(input, ':')
	return input
}

func writeInt(input []byte, val reflect.Value) []byte {
	input = strconv.AppendInt(input, val.Int(), 10)
	return input
}
func writeFloat32(input []byte, val reflect.Value) []byte {
	input = strconv.AppendFloat(input, val.Float(), 'f', -1, 32)
	return input
}
func writeFloat64(input []byte, val reflect.Value) []byte {
	input = strconv.AppendFloat(input, val.Float(), 'f', -1, 64)
	return input
}
func writeBool(input []byte, val reflect.Value) []byte {
	if val.Bool() {
		input = append(input, "true"...)
	} else {
		input = append(input, "false"...)
	}
	return input
}
func writeVal(input []byte, val reflect.Value) []byte {
	input = append(input, val.String()...)
	return input
}
func writeString(input []byte, val string) []byte {
	// 处理换行等特殊符号
	val = strings.Replace(val, "\"", "\\\"", -1)
	val = strings.Replace(val, "\r", "\\r", -1)
	val = strings.Replace(val, "\n", "\\n", -1)
	input = append(input, '"')
	input = append(input, val...)
	input = append(input, '"')
	return input
}

func writeNil(input []byte) []byte {
	input = append(input, "null"...)
	return input
}
func writeUnknown(input []byte) []byte {
	input = append(input, "unknown"...)
	return input
}
func lineBreak(input []byte) []byte {
	input = append(input, '\n')
	return input
}
func comma(input []byte) []byte {
	input = append(input, ',')
	return input
}
func leftBracket(input []byte) []byte {
	input = append(input, '[')
	return input
}
func rightBracket(input []byte) []byte {
	input = append(input, ']')
	return input
}

func WriteObject(input []byte, rv reflect.Value) []byte {
	if rv.Interface() == nil {
		input = writeNil(input)
		return input
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		input = writeString(input, rv.String())
	case reflect.Int, reflect.Int16, reflect.Int64, reflect.Int32, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		input = writeInt(input, rv)
	case reflect.Float32:
		input = writeFloat32(input, rv)
	case reflect.Float64:
		input = writeFloat64(input, rv)
	case reflect.Bool:
		input = writeBool(input, rv)
	case reflect.Slice:
		input = leftBracket(input)
		for i := 0; i < rv.Len(); i++ {
			if i != 0 {
				input = comma(input)
			}
			input = WriteObject(input, rv.Index(i))
		}
		input = rightBracket(input)
	case reflect.Struct:
		input = start(input)
		// do something
		for i := 0; i < rv.NumField(); i++ {
			key := rv.Type().Field(i).Tag.Get("json")
			if key == "" {
				key = rv.Type().Field(i).Name
			}
			subVal := rv.Field(i)
			input = writeKey(input, key)
			input = WriteObject(input, subVal)
			if i < rv.NumField()-1 {
				input = comma(input)
			}
			// lineBreak()
		}
		input = end(input)
	case reflect.Interface:
		rv = rv.Elem()
		input = WriteObject(input, rv)
	default:
		input = writeUnknown(input)
	}
	return input
}
