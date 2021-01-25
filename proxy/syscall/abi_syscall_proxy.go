// How to use conditional compilation with the go build tool:
// https://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool

// +build proxytest

package syscall

import "github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"

type WasmHost interface {
	ProxyLog(logLevel types.LogLevel, buffer *byte, len int) types.Status
	ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) types.Status
	ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) types.Status
	ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status
	ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status
	ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status
	ProxyContinueStream(streamType types.StreamType) types.Status
	ProxyCloseStream(streamType types.StreamType) types.Status
	ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status
	ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status
	ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status
	ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status
	ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status
	ProxySetTickPeriodMilliseconds(period uint32) types.Status
	ProxySetEffectiveContext(contextID uint32) types.Status
	ProxyDone() types.Status
}

var proxyHost WasmHost

func RegisterWasmHost(wasmHost WasmHost) {
	proxyHost = wasmHost
}

type DefaultWasmHost struct {
}

var defaultHost WasmHost = &DefaultWasmHost{}

func (h *DefaultWasmHost) ProxyLog(logLevel types.LogLevel, buffer *byte, len int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyContinueStream(streamType types.StreamType) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyCloseStream(streamType types.StreamType) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetEffectiveContext(contextID uint32) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyDone() types.Status { return types.StatusOK }

func wasmHost() WasmHost {
	if proxyHost == nil {
		return defaultHost
	}
	return proxyHost
}

//export proxy_log
func ProxyLog(logLevel types.LogLevel, buffer *byte, len int) types.Status {
	return wasmHost().ProxyLog(logLevel, buffer, len)
}

//export proxy_get_header_map_value
func ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	return wasmHost().ProxyGetHeaderMapValue(mapType, keyData, keySize, returnValueData, returnValueSize)
}

//export proxy_add_header_map_value
func ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return wasmHost().ProxyAddHeaderMapValue(mapType, keyData, keySize, valueData, valueSize)
}

//export proxy_replace_header_map_value
func ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return wasmHost().ProxyReplaceHeaderMapValue(mapType, keyData, keySize, valueData, valueSize)
}

//export proxy_remove_header_map_value
func ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	return wasmHost().ProxyRemoveHeaderMapValue(mapType, keyData, keySize)
}

//export proxy_get_header_map_pairs
func ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status {
	return wasmHost().ProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
}

//export proxy_set_header_map_pairs
func ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	return wasmHost().ProxySetHeaderMapPairs(mapType, mapData, mapSize)
}

//export proxy_get_buffer_bytes
func ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status {
	return wasmHost().ProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
}

//export proxy_set_buffer_bytes
func ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status {
	return wasmHost().ProxySetBufferBytes(bt, start, maxSize, bufferData, bufferSize)
}

//export proxy_continue_stream
func ProxyContinueStream(streamType types.StreamType) types.Status {
	return wasmHost().ProxyContinueStream(streamType)
}

//export proxy_close_stream
func ProxyCloseStream(streamType types.StreamType) types.Status {
	return wasmHost().ProxyCloseStream(streamType)
}

//export proxy_set_tick_period_milliseconds
func ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	return wasmHost().ProxySetTickPeriodMilliseconds(period)
}

////export proxy_get_current_time_nanoseconds
//func ProxyGetCurrentTimeNanoseconds(returnTime *int64) types.Status {
//	return wasmHost().ProxyGet
//}

//export proxy_set_effective_context
func ProxySetEffectiveContext(contextID uint32) types.Status {
	return wasmHost().ProxySetEffectiveContext(contextID)
}

//export proxy_done
func ProxyDone() types.Status {
	return wasmHost().ProxyDone()
}

//export proxy_get_property
func ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) types.Status {
	return wasmHost().ProxyGetProperty(pathData, pathSize, returnValueData, returnValueSize)
}

//export proxy_set_property
func ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) types.Status {
	return wasmHost().ProxySetProperty(pathData, pathSize, valueData, valueSize)
}
