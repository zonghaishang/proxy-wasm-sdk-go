package crpc

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
	"reflect"
	"strconv"
)

/*

* 0     1     2           4           6           8
* +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+------+-----+-----+-----+-----+
* |magic num  |requestLen             | headLen   |                   header          |        body
* -------------------------------------------Frame---------------------------------------------------
*                                     |------------------------------------Request-------------------
*                                                 |---------header---|------body---------------------



* 0     1     2           4           6           8          10     11     12          14         16
* +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+------+-----+-----+-----+-----+
* |magic num  |requestLen             | headLen   |ver  |  headProperties  |        requestId
* +-----------+-----------+-----------+-----------+-----------+------------+-----------+-----------+
*                           requestId(16 byte)                             | timeout   |
* +-----------+-----------+-----------+-----------+-----------+------------+-----------+-----------+
*             |   tranNum(36 byte)                                                                 |
*                                                                                                  +
*                                                                                                  |
*                                                                                                  +
|                applySysTime (26 byte)                                                            |
* +-----------+-----------+-----------+-----------+-----------+------------+-----------+-----------+

* +-----------+-----------+-----------+-----------+-----------+------------+-----------+-----------+
  | callApp               |TagCnt     |
* +-----------+-----------+-----------+-----------+-----------+------------+-----------+-----------+
TagExtends     {tagName(4 byte) + tagType(1 byte) + tagValueLen(2 byte) + tagValue}      |
* +-----------+-----------+-----------+-----------+-----------+------------+-----------+-----------+
body bytes ... ...
* +-----------+-----------+-----------+-----------+-----------+------------+-----------+-----------+
*/

type crpcCodec struct {
}

func (c *crpcCodec) Encode(ctx context.Context, cmd proxy.Command) (proxy.Buffer, error) {
	if req, ok := cmd.(*Request); ok {
		return encodeRequest(ctx, req)
	}

	// encode response command
	if resp, ok := cmd.(*Response); ok {
		return encodeResponse(ctx, resp)
	}
	proxy.Log.Warnf("[crpcCodec] maybe receive unsupported command type: %v", reflect.TypeOf(cmd))
	return nil, fmt.Errorf("unknow command type: %v", reflect.TypeOf(cmd))
}

func (c *crpcCodec) Decode(ctx context.Context, data proxy.Buffer) (proxy.Command, error) {
	if data.Len() >= RequestHeartBeatHeaderLen {
		//1000,0000 第一位表示请求类型，1=request，0=response
		if data.Bytes()[9]&0x80 == 0x80 {
			if data.Bytes()[9]&0x20 == 0x20 {
				return decodeHeartbeatRequest(ctx, data)
			}
			return decodeRequest(ctx, data)
		} else {
			if data.Bytes()[9]&0x20 == 0x20 {
				return decodeHeartbeatResponse(ctx, data)
			}
			return decodeResponse(ctx, data)
		}
	}

	return nil, nil
}

func encodeRequest(ctx context.Context, request *Request) (proxy.Buffer, error) {
	//fast-path
	if request.rawData != nil {
		if !request.Changed && !isChanged(&request.CommonHeader) && !request.ContentChanged {
			//TODO request.Data.Count(1)
			return request.Data, nil
		}
	}
	// delete service Key
	request.CommonHeader.Del("service")
	// slow-path count
	headerLen, requestLen, frameLen := countRequestMsgHeaderLen(request)
	request.RequestLen = requestLen
	request.HeaderLen = headerLen
	request.TagCnt = byte(len(request.ToMap()))
	buf := proxy.NewBuffer(frameLen)
	buf.WriteByte(ProtocolFirstByte)
	buf.WriteByte(ProtocolSecondByte)
	buf.WriteUint32(request.RequestLen)
	buf.WriteUint16(request.HeaderLen)
	buf.WriteByte(request.Version)
	buf.Write(request.HeaderProperties)
	uid, _ := Parse(request.RequestId)
	buf.Write(uid[:])
	if !request.Heartbeat {
		buf.WriteUint32(request.Timeout)
		buf.WriteString(request.SourceApp)
		buf.WriteString(request.TranNum)
		buf.WriteString(request.ApplySysTime)
		buf.WriteString(request.CallApp)
		buf.WriteByte(request.TagCnt)
	}
	//removeRequestHeader(request)
	encodeTags(buf, &request.CommonHeader)
	if request.Body != nil {
		buf.Write(request.Body.Bytes())
	}

	return buf, nil
}

func encodeTags(buf proxy.Buffer, header *proxy.CommonHeader) {
	for k, v := range header.ToMap() {
		keyLen := byte(len(k))
		buf.WriteByte(keyLen)
		buf.WriteString(k)
		buf.WriteByte(0x01)
		buf.WriteUint16(uint16(len(v)))
		buf.WriteString(v)
	}
}

func countRequestMsgHeaderLen(request *Request) (uint16, uint32, int) {
	// heartbeat return
	if request.Heartbeat {
		return request.HeaderLen, request.RequestLen, int(request.RequestLen) + HeaderBeginLen
	}
	// 计算
	headerLen := RequestTagStartIndex
	headerLen = headerLen + getHeaderEncodeLength(&request.CommonHeader)
	frameLen := headerLen
	if request.Body != nil {
		frameLen = headerLen + len(request.Body.Bytes())
	}
	return uint16(headerLen - HeaderBeginLen - 2), uint32(frameLen - HeaderBeginLen), frameLen
}

func getHeaderEncodeLength(header *proxy.CommonHeader) (size int) {
	h := header.ToMap()
	for k, v := range h {
		size += 4 + len(k) + len(v)
	}
	return
}

func isChanged(header proxy.Header) bool {
	if _, ok := header.Get(SOFARPC_ROUTER_RULE_METADATA_CRPC_TARGET_HOST); ok {
		return true
	}
	if _, ok := header.Get(SOFARPC_ROUTER_RULE_METADATA_CRPC_TARGET_ZONE); ok {
		return true
	}
	return false
}

func encodeResponse(ctx context.Context, response *Response) (proxy.Buffer, error) {
	if response.rawData != nil {
		if !response.Changed && isChanged(&response.CommonHeader) && !response.ContentChanged {
			// hack: increase the buffer count to avoid premature recycle
			//TODO response.Data.Count(1)
			return response.Data, nil
		}
	}
	response.CommonHeader.Del("service")
	headerLen, responseLen, frameLen := countResponseMsgHeaderLen(response)
	response.ResponseLen = responseLen
	response.HeaderLen = headerLen
	response.TagCnt = byte(len(response.ToMap()))
	buf := proxy.NewBuffer(frameLen)

	buf.WriteByte(ProtocolFirstByte)
	buf.WriteByte(ProtocolSecondByte)
	buf.WriteUint32(response.ResponseLen)
	buf.WriteUint16(response.HeaderLen)
	buf.WriteByte(response.Version)
	buf.Write(response.HeaderProperties)
	uid, _ := Parse(response.RequestId)
	buf.Write(uid[:])
	if !response.Heartbeat {
		buf.WriteString(response.TranNum)
		buf.WriteString(response.RpcRespCode)
		buf.WriteString(response.AppRespCode)
		buf.WriteByte(response.TagCnt)
	}

	//removeResponseHeader(response)
	encodeTags(buf, &response.CommonHeader)
	if response.Body != nil {
		buf.Write(response.Body.Bytes())
	}

	return buf, nil
}

func countResponseMsgHeaderLen(response *Response) (uint16, uint32, int) {
	// heartbeat return
	if response.Heartbeat {
		return response.HeaderLen, response.ResponseLen, int(response.ResponseLen) + HeaderBeginLen
	}
	// 计算
	headerLen := ResponseTagStartIndex
	headerLen = headerLen + getHeaderEncodeLength(&response.CommonHeader)
	frameLen := headerLen
	if response.Body != nil {
		frameLen = headerLen + len(response.Body.Bytes())
	}
	return uint16(headerLen - HeaderBeginLen - 2), uint32(frameLen - HeaderBeginLen), frameLen
}

func decodeRequest(ctx context.Context, data proxy.Buffer) (proxy.Command, error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. least bytes to decode header is RequestHeaderLen(100)
	if bytesLen < RequestHeaderLen {
		return nil, nil
	}

	//2. least bytes to decode whole frame
	requestLen := binary.BigEndian.Uint32(bytes[2:6])
	headerLen := binary.BigEndian.Uint16(bytes[6:8])
	bodyLen := requestLen - uint32(headerLen) - 2

	frameLen := int(requestLen) + HeaderBeginLen
	if bytesLen < frameLen {
		return nil, nil
	}

	data.Drain(frameLen)

	request := &Request{}

	//3. decode header
	request.RequestHeader = RequestHeader{
		RequestLen:       requestLen,
		HeaderLen:        headerLen,
		Version:          bytes[8],
		HeaderProperties: copyByte(bytes[9:12]),
		Timeout:          binary.BigEndian.Uint32(bytes[28:32]),
		SourceApp:        string(bytes[32:36]),
		TranNum:          string(bytes[36:72]),
		ApplySysTime:     string(bytes[72:98]),
		CallApp:          string(bytes[98:102]),
		TagCnt:           bytes[102],
	}

	request.RequestId = getUUID(bytes[12:28])
	if ctx.Value(types.ContextKeyListenerType) == EGRESS {
		request.isEgress = true
	} else {
		request.isEgress = false
	}

	request.OneWay = request.HeaderProperties[0]&0x64 == 0x64
	request.Heartbeat = false

	request.Data = proxy.NewBuffer(frameLen)

	//4. set timeout to notify proxy
	//TODO variable.SetVariableValue(ctx, VarProxyGlobalTimeout, strconv.Itoa(int(request.Timeout)))

	//5. copy data for io multiplexing
	request.Data.Write(bytes[:frameLen])
	request.rawData = request.Data.Bytes()

	//6. process wrappers: Header, Tag, Body
	tagIndex := RequestTagStartIndex

	//6.1 process tags
	request.rawTags = request.rawData[tagIndex : tagIndex+(int(headerLen)-RequestHeaderLenBeforeTag)]
	err := decodeTag(int(request.TagCnt), request.rawTags, &request.CommonHeader)

	//6.2 process body
	request.rawBody = request.rawData[frameLen-int(bodyLen):]

	request.Body = proxy.NewBuffer(len(request.rawBody))
	request.Body.Write(request.rawBody)

	serviceName, ok := request.Get(SERVICE_NAME_KEY)
	if ok {
		request.CommonHeader.Set("service", serviceName)
	}

	SetRequestHeaderValue(request)
	return request, err
}

func decodeHeartbeatRequest(ctx context.Context, data proxy.Buffer) (proxy.Command, error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. least bytes to decode header is RequestHeaderLen(100)
	if bytesLen < RequestHeartBeatHeaderLen {
		return nil, nil
	}

	//2. least bytes to decode whole frame
	requestLen := binary.BigEndian.Uint32(bytes[2:6])
	headerLen := binary.BigEndian.Uint16(bytes[6:8])
	bodyLen := requestLen - uint32(headerLen) - 2

	frameLen := int(requestLen) + HeaderBeginLen
	if bytesLen < frameLen {
		return nil, nil
	}

	data.Drain(frameLen)

	request := &Request{}

	//3. decode header
	request.RequestHeader = RequestHeader{
		RequestLen:       requestLen,
		HeaderLen:        headerLen,
		Version:          bytes[8],
		HeaderProperties: copyByte(bytes[9:12]),
	}

	request.RequestId = getUUID(bytes[12:28])

	if ctx.Value(types.ContextKeyListenerType) == EGRESS {
		request.isEgress = true
	} else {
		request.isEgress = false
	}

	request.OneWay = request.HeaderProperties[0]&0x64 == 0x64
	request.Heartbeat = true

	request.Data = proxy.NewBuffer(frameLen)

	//5. copy data for io multiplexing
	request.Data.Write(bytes[:frameLen])
	request.rawData = request.Data.Bytes()

	//6.2 process body
	request.rawBody = request.rawData[frameLen-int(bodyLen):]

	request.Body = proxy.NewBuffer(len(request.rawBody))
	request.Body.Write(request.rawBody)
	request.CommonHeader.Set("service", "heartbeatReq")

	return request, nil
}

func SetRequestHeaderValue(request *Request) {
	// service id
	if _, ok := request.Get(GOVERN_SERVICE_KEY); !ok {
		var serviceId string
		value, ok := request.Get(SERVICE_NAME_KEY)
		if ok {
			serviceId = value
		} else {
			value, ok := request.Get(SERVICE_NAME_KEY)
			if ok {
				serviceId = value
			}
		}
		//service version
		value, ok = request.Get(SERVICE_VERSION_KEY)
		if ok && value != "" {
			serviceId = serviceId + ":" + value
		} else {
			serviceId = serviceId + ":1.0.0"
		}
		//group Id
		value, ok = request.Get(GROUP_ID_KEY)
		if ok && value != "" {
			serviceId = serviceId + ":" + value
		} else {
			serviceId = serviceId + ":default"
		}
		// put service
		if serviceId != "" {
			request.Set(GOVERN_SERVICE_KEY, serviceId+"@crpc")
		}
	}
	// service type : crpc
	if _, ok := request.Get(GOVERN_SERVICE_KEY); !ok {
		request.Set(GOVERN_SERVICE_TYPE_KEY, PROTOCOL_NAME)
	}

	// method name
	if _, ok := request.Get(GOVERN_METHOD_KEY); !ok {
		value, ok := request.Get(SERVICE_METHOD_NAME_KEY)
		if ok {
			request.Set(GOVERN_METHOD_KEY, value)
		}
	}

	// timeout
	if _, ok := request.Get(GOVERN_TIMEOUT_KEY); !ok {
		request.Set(GOVERN_TIMEOUT_KEY, strconv.Itoa(int(request.Timeout)))
	}

	// target app name
	if _, ok := request.Get(GOVERN_TARGET_APP_KEY); !ok {
		value, ok := request.Get(TARGET_APP_NAME_KEY)
		if ok {
			request.Set(GOVERN_TARGET_APP_KEY, value)
		}
	}

	//source AppName
	if _, ok := request.Get(GOVERN_SOURCE_APP_KEY); !ok {
		request.Set(GOVERN_SOURCE_APP_KEY, request.SourceApp)
	}

	// request Id
	if _, ok := request.Get(GOVERN_REQUEST_ID); !ok {
		request.Set(GOVERN_REQUEST_ID, request.RequestId)
	}

	// tracer Id
	if _, ok := request.Get(GOVERN_TRACE_ID); !ok {
		value, ok := request.Get(TRACE_ID_KEY)
		if ok {
			request.Set(GOVERN_TRACE_ID, value)
		} else {
			//TODO request.Set(GOVERN_TRACE_ID, trace.IdGen().GenerateTraceId())
		}
	}

	// span Id
	if _, ok := request.Get(GOVERN_TRACE_SPAN_ID); !ok {
		value, ok := request.Get(SPAN_ID_KEY)
		if ok {
			request.Set(GOVERN_TRACE_SPAN_ID, value)
		}
	}

	// parent span Id
	if _, ok := request.Get(GOVERN_TRACE_PARENT_SPAN_ID); !ok {
		value, ok := request.Get(PARENT_SPAN_ID_KEY)
		if ok {
			request.Set(GOVERN_TRACE_PARENT_SPAN_ID, value)
		}
	}

	// sampled
	if _, ok := request.Get(GOVERN_TRACE_SAMPLED); !ok {
		value, ok := request.Get(SAMPLED_KEY)
		if ok {
			request.Set(GOVERN_TRACE_SAMPLED, value)
		}
	}

	// flags
	if _, ok := request.Get(GOVERN_TRACE_flags); !ok {
		value, ok := request.Get(FLAGS_KEY)
		if ok {
			request.Set(GOVERN_TRACE_flags, value)
		}
	}

	// tran num
	request.CommonHeader.Set(TRAN_NUM, request.TranNum)
}

func decodeResponse(ctx context.Context, data proxy.Buffer) (cmd proxy.Command, err error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. least bytes to decode header is ResponseHeaderLen(89)
	if bytesLen < ResponseTagStartIndex {
		return
	}

	// 2. least bytes to decode whole frame
	responseLen := binary.BigEndian.Uint32(bytes[2:6])
	headerLen := binary.BigEndian.Uint16(bytes[6:8])
	bodyLen := responseLen - uint32(headerLen) - 2

	frameLen := int(responseLen) + HeaderBeginLen
	if bytesLen < frameLen {
		return
	}

	data.Drain(frameLen)

	response := &Response{}

	response.ResponseHeader = ResponseHeader{
		ResponseLen:      responseLen,
		HeaderLen:        headerLen,
		Version:          bytes[8],
		HeaderProperties: copyByte(bytes[9:12]),
		TranNum:          string(bytes[28:64]),
		RpcRespCode:      string(bytes[64:71]),
		AppRespCode:      string(bytes[71:78]),
		TagCnt:           bytes[78],
	}

	response.Heartbeat = false

	response.RequestId = getUUID(bytes[12:28])
	response.Data = proxy.NewBuffer(frameLen)
	response.Data.Write(bytes[:frameLen])
	response.rawData = response.Data.Bytes()

	//6.1 process tags
	response.rawTags = response.rawData[ResponseTagStartIndex : ResponseTagStartIndex+(int(headerLen)-ResponseHeaderLenBeforeTag)]
	err = decodeTag(int(response.TagCnt), response.rawTags, &response.CommonHeader)

	//6.2 process body
	response.rawBody = response.rawData[frameLen-int(bodyLen):]
	response.Body = proxy.NewBuffer(len(response.rawBody))
	response.Body.Write(response.rawBody)
	response.CommonHeader.Set("service", response.RequestId)

	setResponseHeaderValue(response)
	return response, err
}

func setResponseHeaderValue(response *Response) {
	response.Set(GOVERN_REQUEST_ID, response.RequestId)
	responseCode, _ := MappingCrpcCode2HttpCode(response)
	response.Set(GOVERN_RESP_CODE_KEY, strconv.Itoa(responseCode))
	response.Set(CRPC_TRACER_HEADER_RESPONSE_CRPC_CODE, response.RpcRespCode)
	response.Set(CRPC_TRACER_HEADER_RESPONSE_BIZ_CODE, response.AppRespCode)

}

func decodeHeartbeatResponse(ctx context.Context, data proxy.Buffer) (proxy.Command, error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. least bytes to decode header is RequestHeaderLen(100)
	if bytesLen < RequestHeartBeatHeaderLen {
		return nil, nil
	}

	//2. least bytes to decode whole frame
	responseLen := binary.BigEndian.Uint32(bytes[2:6])
	headerLen := binary.BigEndian.Uint16(bytes[6:8])
	bodyLen := responseLen - uint32(headerLen) - 2

	frameLen := int(responseLen) + HeaderBeginLen
	if bytesLen < frameLen {
		return nil, nil
	}

	data.Drain(frameLen)

	response := &Response{}

	//3. decode header
	response.ResponseHeader = ResponseHeader{
		ResponseLen:      responseLen,
		HeaderLen:        headerLen,
		Version:          bytes[8],
		HeaderProperties: copyByte(bytes[9:12]),
	}

	//response.RequestId = string(bytes[12:28])
	response.RequestId = getUUID(bytes[12:28])
	response.Heartbeat = true

	response.Data = proxy.NewBuffer(frameLen)

	//5. copy data for io multiplexing
	response.Data.Write(bytes[:frameLen])
	response.rawData = response.Data.Bytes()

	//6.2 process body
	response.rawBody = response.rawData[frameLen-int(bodyLen):]

	response.Body = proxy.NewBuffer(len(response.rawBody))
	response.Body.Write(response.rawBody)
	response.CommonHeader.Set("service", "heartbeatResp")

	return response, nil
}

func copyByte(src []byte) []byte {
	m := len(src)
	re := make([]byte, m)
	copy(re, src)
	return re
}

func decodeTag(tagCnt int, bytes []byte, h *proxy.CommonHeader) (err error) {
	index := 0
	for i := 0; i < tagCnt; i++ {
		keyLen := int(bytes[index])
		index++
		key := string(bytes[index : index+keyLen])
		index = index + keyLen
		_ = bytes[index]
		index++
		valueLen := int(binary.BigEndian.Uint16(bytes[index : index+2]))
		index = index + 2
		value := string(bytes[index : index+valueLen])
		h.Set(key, value)
		index = index + valueLen
	}
	return nil
}
