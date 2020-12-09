package spec

import (
	"reflect"
	"unsafe"
)

// parseBytes parse string to byte pointer
func parseBytePtr(message string) *byte {
	if len(message) == 0 {
		buffer := make([]byte, 1)
		return &buffer[0]
	}

	buffer := *(*[]byte)(unsafe.Pointer(&message))
	return &buffer[0]
}

// parseString parse byte pointer to string
func parseString(buf *byte, len int) string {
	if len <= 0 || buf == nil {
		return ""
	}

	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buf)),
		Len:  len,
		Cap:  len,
	}))
}

func parseSliceString(buf []byte) string {
	if len(buf) <= 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&buf))
}

// parseByteSlice parse byte pointer to byte slice
func parseByteSlice(buf *byte, len int) []byte {
	if len <= 0 || buf == nil {
		return []byte{}
	}

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buf)),
		Len:  len,
		Cap:  len,
	}))
}
