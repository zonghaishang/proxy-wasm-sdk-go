package internal

import (
	"reflect"
	"unsafe"
)

// parseBytes parse string to byte pointer
func parseBytePtr(message string) *byte {
	if len(message) == 0 {
		return nil
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

// parseByteSlice parse byte pointer to byte slice
func parseByteSlice(buf *byte, len int) []byte {
	if len <= 0 || buf == nil {
		return nil
	}

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buf)),
		Len:  len,
		Cap:  len,
	}))
}
