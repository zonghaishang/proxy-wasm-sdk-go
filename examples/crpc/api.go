package crpc

import "github.com/zonghaishang/proxy-wasm-sdk-go/proxy"

func NewRequest(requestId string, headers proxy.Header, data proxy.Buffer) *Request {
	request := &Request{
		RequestHeader: RequestHeader{
			MagicNum:         []byte{ProtocolFirstByte, ProtocolSecondByte},
			RequestLen:       uint32(RequestHeaderLen),
			HeaderLen:        uint16(RequestHeaderLen),
			Version:          0x01,
			HeaderProperties: []byte{0xc4, 0x0, 0x0}, //1100 0000
		},
	}
	if headers != nil {
		headers.Range(func(key, value string) bool {
			request.Set(key, value)
			return true
		})
	}
	if data != nil {
		//request.Data = data
		request.Body = data
	}
	return request
}

func NewResponse(requestId string, statusCode string, headers proxy.Header, data proxy.Buffer) *Response {
	response := &Response{
		ResponseHeader: ResponseHeader{
			MagicNum:         []byte{ProtocolFirstByte, ProtocolSecondByte},
			ResponseLen:      uint32(HeartBeatRequestHeaderLen),
			HeaderLen:        uint16(HeartBeatHeaderLen),
			Version:          0x01,
			HeaderProperties: []byte{0x0, 0x0, 0x0}, //0000 0000
			RequestId:        requestId,
			RpcRespCode:      statusCode,
		},
	}
	// set headers
	if headers != nil {
		headers.Range(func(key, value string) bool {
			response.Set(key, value)
			return true
		})
	}

	// set content
	if data != nil {
		response.Data = data
	}
	return response
}
