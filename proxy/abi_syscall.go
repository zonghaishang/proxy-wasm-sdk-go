// How to use conditional compilation with the go build tool:
// https://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool

// +build !proxytest

package proxy

import "github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"

//export proxy_log
func ABI_ProxyLog(logLevel types.LogLevel, buffer *byte, len int) types.Status

//export proxy_get_header_map_value
func ABI_ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status

//export proxy_add_header_map_value
func ABI_ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status

//export proxy_replace_header_map_value
func ABI_ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status

//export proxy_remove_header_map_value
func ABI_ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status

//export proxy_get_header_map_pairs
func ABI_ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status

//export proxy_set_header_map_pairs
func ABI_ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status

//export proxy_get_buffer_bytes
func ABI_ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status

//export proxy_set_buffer_bytes
func ABI_ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status

//export proxy_continue_stream
func ABI_ProxyContinueStream(streamType types.StreamType) types.Status

//export proxy_close_stream
func ABI_ProxyCloseStream(streamType types.StreamType) types.Status

//export proxy_set_tick_period_milliseconds
func ABI_ProxySetTickPeriodMilliseconds(period uint32) types.Status

//export proxy_get_current_time_nanoseconds
func ABI_ProxyGetCurrentTimeNanoseconds(returnTime *int64) types.Status

//export proxy_set_effective_context
func ABI_ProxySetEffectiveContext(contextID uint32) types.Status

//export proxy_done
func ABI_ProxyDone() types.Status

//export proxy_get_property
func ABI_ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) types.Status

//export proxy_set_property
func ABI_ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) types.Status
