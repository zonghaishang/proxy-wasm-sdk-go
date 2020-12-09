package internal

import "github.com/zonghaishang/proxy-wasm-sdk-go/internal/types"

//export proxy_get_buffer
func proxyGetBuffer(bt types.BufferType, offset int, maxSize int, buf **byte, len *int) types.Status {
	return types.StatusOK
}

//export proxy_set_buffer
func proxySetBuffer(bt types.BufferType, offset int, maxSize int, buf *byte, len int) types.Status {
	return types.StatusOK
}

//export proxy_get_map_values
func proxyGetMapValues(mt types.MapType, keys *byte, keyCount int, pairs **byte, pairCount *int) types.Status {
	return types.StatusOK
}

//export proxy_set_map_values
func proxySetMapValues(mt types.MapType, removeKeys *byte, removeKeyCount int, pairs **byte, pairCount int) types.Status {
	return types.StatusOK
}

//export proxy_log
func proxyLog(logLevel types.LogLevel, buf *byte, len int) types.Status {
	return types.StatusOK
}
