package spec

import "github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"

func GetPluginConfiguration(size int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypePluginConfiguration, 0, size)
	return buf, types.StatusToError(status)
}

func GetVMConfiguration(size int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypeVMConfiguration, 0, size)
	return buf, types.StatusToError(status)
}

func SetTickPeriodMilliSeconds(millSec uint32) error {
	return types.StatusToError(ProxySetTickPeriodMilliseconds(millSec))
}

func GetDownStreamData(start, maxSize int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypeDownstreamData, start, maxSize)
	return buf, types.StatusToError(status)
}

func GetUpstreamData(start, maxSize int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypeUpstreamData, start, maxSize)
	return buf, types.StatusToError(status)
}

func GetHttpRequestHeaders() (map[string]string, error) {
	headers, status := getMap(types.MapTypeHttpRequestHeaders)
	return headers, types.StatusToError(status)
}

func SetHttpRequestHeaders(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestHeaders, headers))
}

func GetHttpRequestHeader(key string) (string, error) {
	header, status := getMapValue(types.MapTypeHttpRequestHeaders, key)
	return header, types.StatusToError(status)
}

func RemoveHttpRequestHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestHeaders, key))
}

func SetHttpRequestHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

func AddHttpRequestHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

func GetHttpRequestBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpRequestBody, start, maxSize)
	return ret, types.StatusToError(st)
}

func SetHttpRequestBody(body []byte) error {
	var buff *byte
	if len(body) != 0 {
		buff = &body[0]
	}
	status := ProxySetBufferBytes(types.BufferTypeHttpRequestBody, 0, len(body), buff, len(body))
	return types.StatusToError(status)
}

func GetHttpRequestTrailers() (map[string]string, error) {
	trailers, status := getMap(types.MapTypeHttpRequestTrailers)
	return trailers, types.StatusToError(status)
}

func SetHttpRequestTrailers(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestTrailers, headers))
}

func GetHttpRequestTrailer(key string) (string, error) {
	trailer, status := getMapValue(types.MapTypeHttpRequestTrailers, key)
	return trailer, types.StatusToError(status)
}

func RemoveHttpRequestTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestTrailers, key))
}

func SetHttpRequestTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

func AddHttpRequestTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

func ResumeHttpRequest() error {
	return types.StatusToError(ProxyContinueStream(types.StreamTypeRequest))
}

func GetHttpResponseHeaders() (map[string]string, error) {
	headers, status := getMap(types.MapTypeHttpResponseHeaders)
	return headers, types.StatusToError(status)
}

func SetHttpResponseHeaders(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseHeaders, headers))
}

func GetHttpResponseHeader(key string) (string, error) {
	header, status := getMapValue(types.MapTypeHttpResponseHeaders, key)
	return header, types.StatusToError(status)
}

func RemoveHttpResponseHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseHeaders, key))
}

func SetHttpResponseHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

func AddHttpResponseHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

func GetHttpResponseBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpResponseBody, start, maxSize)
	return ret, types.StatusToError(st)
}

func SetHttpResponseBody(body []byte) error {
	var buf *byte
	if len(body) != 0 {
		buf = &body[0]
	}
	st := ProxySetBufferBytes(types.BufferTypeHttpResponseBody, 0, len(body), buf, len(body))
	return types.StatusToError(st)
}

func GetHttpResponseTrailers() (map[string]string, error) {
	trailers, status := getMap(types.MapTypeHttpResponseTrailers)
	return trailers, types.StatusToError(status)
}

func SetHttpResponseTrailers(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseTrailers, headers))
}

func GetHttpResponseTrailer(key string) (string, error) {
	trailer, status := getMapValue(types.MapTypeHttpResponseTrailers, key)
	return trailer, types.StatusToError(status)
}

func RemoveHttpResponseTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseTrailers, key))
}

func SetHttpResponseTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

func AddHttpResponseTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

func ResumeHttpResponse() error {
	return types.StatusToError(ProxyContinueStream(types.StreamTypeResponse))
}

func GetProperty(path []string) ([]byte, error) {
	var ret *byte
	var retSize int
	raw := EncodePropertyPath(path)

	err := types.StatusToError(ProxyGetProperty(&raw[0], len(raw), &ret, &retSize))
	if err != nil {
		return nil, err
	}

	return parseByteSlice(ret, retSize), nil
}

func SetProperty(path string, data []byte) error {
	return types.StatusToError(ProxySetProperty(
		stringBytePtr(path), len(path), &data[0], len(data),
	))
}

func setMap(mapType types.MapType, headers map[string]string) types.Status {
	encodedBytes := EncodeMap(headers)
	hp := &encodedBytes[0]
	hl := len(encodedBytes)
	return ProxySetHeaderMapPairs(mapType, hp, hl)
}

func getMapValue(mapType types.MapType, key string) (string, types.Status) {
	var rvs int
	var raw *byte
	if st := ProxyGetHeaderMapValue(mapType, stringBytePtr(key), len(key), &raw, &rvs); st != types.StatusOK {
		return "", st
	}

	ret := parseString(raw, rvs)
	return ret, types.StatusOK
}

func removeMapValue(mapType types.MapType, key string) types.Status {
	return ProxyRemoveHeaderMapValue(mapType, stringBytePtr(key), len(key))
}

func setMapValue(mapType types.MapType, key, value string) types.Status {
	return ProxyReplaceHeaderMapValue(mapType, stringBytePtr(key), len(key), stringBytePtr(value), len(value))
}

func addMapValue(mapType types.MapType, key, value string) types.Status {
	return ProxyAddHeaderMapValue(mapType, stringBytePtr(key), len(key), stringBytePtr(value), len(value))
}

func getMap(mapType types.MapType) (map[string]string, types.Status) {
	var rvs int
	var raw *byte

	status := ProxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if status != types.StatusOK {
		return nil, status
	}

	bs := parseByteSlice(raw, rvs)
	return DecodeMap(bs), types.StatusOK
}

func getBuffer(bufType types.BufferType, start, maxSize int) ([]byte, types.Status) {
	var buffer *byte
	var len int
	switch status := ProxyGetBufferBytes(bufType, start, maxSize, &buffer, &len); status {
	case types.StatusOK:
		// is this correct handling...?
		if buffer == nil {
			return nil, types.StatusNotFound
		}
		return parseByteSlice(buffer, len), status
	default:
		return nil, status
	}
}
