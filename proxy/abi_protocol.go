package proxy

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
)

//export proxy_decode_buffer_bytes
func proxyDecodeBufferBytes(contextID uint32, bufferData **byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		panic("invalid context on proxy_decode_buffer_bytes")
	}
	this.setActiveContextID(contextID)

	if len <= 0 {
		// should never be happen
		return types.StatusEmpty
	}

	// convert data into an array of bytes to be parsed
	data := parseByteSlice(*bufferData, len)
	buffer := Allocate(data)
	// call user extension implementation
	cmd, err := ctx.Codec().Decode(buffer)

	// need more data
	if cmd == nil {
		return types.StatusNeedMoreData
	}

	// we check decode is ok
	if err != nil {
		log.Fatalf("failed to decode buffer by protocol %s, context id %v", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}

	// bufferData format:
	// encoded header map | Flag | Id | (Timeout|Status) | drain length | raw bytes

	return types.StatusOK
}

//export proxy_encode_buffer_bytes
func pProxyEncodeBufferBytes(contextID uint32, bufferData **byte, bufferSize int) types.Status {
	return types.StatusOK
}

//export proxy_keepalive_buffer_bytes
func proxyKeepAliveBufferBytes(contextID uint32, id uint64) types.Status {
	return types.StatusOK
}

//export proxy_reply_keepalive_buffer_bytes
func proxyReplyKeepAliveBufferBytes(contextID uint32, bufferData *byte, bufferSize int) types.Status {
	return types.StatusOK
}

//export proxy_hijack_buffer_bytes
func proxyHijackBufferBytes(contextID uint32, statusCode uint32, bufferData *byte, bufferSize int) types.Status {
	return types.StatusOK
}
