package crpc

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"strconv"
)

type crpcProtocol struct {
	crpcCodec
	proxy.DefaultOptions
}

func (c crpcProtocol) Name() string {
	return PROTOCOL_NAME
}

func (c crpcProtocol) Codec() proxy.Codec {
	return &c.crpcCodec
}

func (c crpcProtocol) KeepAlive(requestId uint64) proxy.Request {
	request := &Request{
		RequestHeader: RequestHeader{
			MagicNum:         []byte{ProtocolFirstByte, ProtocolSecondByte},
			RequestLen:       uint32(HeartBeatRequestHeaderLen),
			HeaderLen:        uint16(HeartBeatHeaderLen),
			Version:          0x01,
			HeaderProperties: []byte{0xe0, 0x0, 0x0}, //0110 0000
		},
	}
	uid, err := NewUUID()
	if err != nil {
		return nil
	}
	request.RequestId = uid.String()
	//heartBeatRequestIdCache.Set(request.RequestId, requestId, DefaultExpiration)
	request.CommonHeader.Set("service", "test")
	request.Heartbeat = true
	proxy.Log.Debugf("[heartbeat trigger] requestId " + request.RequestId)
	return request
}

func (c crpcProtocol) ReplyKeepAlive(request proxy.Request) proxy.Response {
	crpcRequest, ok := request.(*Request)
	if !ok {
		return nil
	}

	response := &Response{
		ResponseHeader: ResponseHeader{
			MagicNum:         []byte{ProtocolFirstByte, ProtocolSecondByte},
			ResponseLen:      uint32(HeartBeatRequestHeaderLen),
			HeaderLen:        uint16(HeartBeatHeaderLen),
			Version:          0x01,
			HeaderProperties: []byte{0x32, 0x0, 0x0}, //0010 0000
			RequestId:        crpcRequest.RequestId,
		},
	}
	response.Heartbeat = true
	response.CommonHeader.Set("service", "heartbeat")
	proxy.Log.Debugf("[heartbeat reply] requestId " + response.RequestId)
	return response
}

func (c crpcProtocol) Hijack(request proxy.Request, statusCode uint32) proxy.Response {
	crpcRequest, ok := request.(*Request)
	if !ok {
		return nil
	}
	governCode, ok := GetGovernValue(nil, request.GetHeader(), GOVERN_HIJACK_CODE_KEY)
	if ok {
		code, err := strconv.ParseUint(governCode, 10, 32)
		if err == nil {
			statusCode = uint32(code)
		}
	}
	response := &Response{
		ResponseHeader: ResponseHeader{
			MagicNum:         []byte{ProtocolFirstByte, ProtocolSecondByte},
			ResponseLen:      uint32(66),
			HeaderLen:        uint16(64),
			Version:          0x01,
			HeaderProperties: []byte{0x0, 0x0, 0x0}, //0010 0000
			RequestId:        crpcRequest.RequestId,
			TranNum:          crpcRequest.TranNum,
			RpcRespCode:      MappingHttpCode2CrpcCode(statusCode),
			AppRespCode:      "AAAAAAA",
			TagCnt:           1,
		},
	}
	proxy.Log.Infof("[Hijack] crp hijack. statusCode = %s, RequestId = %s, RpcRespCode = %s",
		strconv.Itoa(int(statusCode)), response.RequestId, response.RpcRespCode)
	response.CommonHeader.Set("service", "hiJack")

	return response
}

func NewCrpcProtocol() proxy.Protocol {
	return &crpcProtocol{}
}

func MappingHttpCode2CrpcCode(httpStatusCode uint32) string {
	switch httpStatusCode {
	case StatusOK:
		return CRPC_SUCCESS
	case RouteHiJackCode:
		return MRPC_ROUTE_ERROR
	case AuthHiJackCode:
		return MRPC_AUTH_ERROR
	case CircuitBreakHiJackCode:
		return MRPC_CIRCUIT_ERROR
	case LimitExceededHiJackCode:
		return MRPC_LIMIT_ERROR
	case FaultInjectHiJackCode:
		return MRPC_FAULT_RULE_ERROR
	case DownGradeHiJackCode:
		return MRPC_DOWNGROUD_ERROR
	case proxy.NoHealthUpstreamCode:
		return MRPC_ROUTE_ERROR
	case proxy.UpstreamOverFlowCode:
		return MRPC_RESPONSE_ERROR_CRPC002
	case proxy.CodecExceptionCode:
		return MRPC_REQUEST_ERROR_CRPC001
	case proxy.DeserializeExceptionCode:
		return MRPC_REQUEST_ERROR_CRPC001
	case proxy.TimeoutExceptionCode:
		return MRPC_TIMEOUT_CRPC003
	default:
		return MRPC_UNKNOWN_ERROR
	}
}

func MappingCrpcCode2HttpCode(response *Response) (int, error) {
	switch response.RpcRespCode {
	case CRPC_SUCCESS:
		return StatusOK, nil
	case CRPC_TIMEOUT:
		return StatusGatewayTimeout, nil
	case CRPC_RPC_REQUEST_ERROR:
		return StatusBadRequest, nil
	case CRPC_ERROR, CRPC_RPC_RESPONSE_ERROR:
		return StatusInternalServerError, nil
	default:
		return StatusInternalServerError, nil
	}
}
