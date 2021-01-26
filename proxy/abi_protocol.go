package proxy

import (
	"context"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
)

//export proxy_decode_buffer_bytes
func proxyDecodeBufferBytes(contextID uint32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode buffer by protocol %s, contextId %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}
	this.setActiveContextID(contextID)

	if len <= 0 {
		// should never be happen
		return types.StatusEmpty
	}

	// convert data into an array of bytes to be parsed
	data := parseByteSlice(bufferData, len)
	buffer := Allocate(data)
	// call user extension implementation
	cmd, err := ctx.Codec().Decode(context.TODO(), buffer)
	if err != nil {
		log.Fatalf("failed to decode buffer by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	// need more data
	if cmd == nil {
		return types.StatusNeedMoreData
	}

	if buffer.Pos() == 0 {
		// When decoding is complete, the contents of the buffer should be read
		log.Errorf("the contents of the buffer should be read by protocol %s, contextId %v, buffer pos %v", ctx.Name(), contextID, buffer.Pos())
		return types.StatusInternalFailure
	}

	ctx.(attribute).set(types.AttributeKeyDecodeCommand, cmd)

	decode := decodeCommandBuffer(cmd, buffer.Pos())
	// report encode data
	err = setDecodeBuffer(decode.Bytes())
	if err != nil {
		log.Errorf("failed to report decode buffer by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	return types.StatusOK
}

//export proxy_encode_buffer_bytes
func proxyEncodeBufferBytes(contextID uint32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode buffer by protocol %s, context replacedId %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}
	this.setActiveContextID(contextID)

	if len <= 0 {
		// should never be happen
		return types.StatusEmpty
	}

	// convert data into an array of dataBytes to be parsed
	data := parseByteSlice(bufferData, len)
	buffer := Allocate(data)

	// bufferData format:
	// encoded header map | Flag | replaceId, id | (Timeout|GetStatus) | drain length | raw dataBytes
	headerBytes, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to read decode buffer header map, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	headers := &CommonHeader{}
	// encoded header map
	if headerBytes > 0 {
		DecodeHeader(data[4:4+headerBytes], headers)
	}

	flag, err := buffer.ReadByte()
	if err != nil {
		log.Errorf("failed to decode buffer flag, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	attr := ctx.(attribute)

	// find context cmd
	cachedCmd := attr.attr(types.AttributeKeyDecodeCommand)
	if cachedCmd == nil {
		// is heartbeat 、keep-alive or hijack ?
		cachedCmd = attr.attr(types.AttributeKeyEncodeCommand)
	}

	if cachedCmd == nil {
		log.Errorf("failed to find cached command, maybe a bug occurred, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	// Multiplexing ID: This is equivalent to the stream ID
	replacedId, err := buffer.ReadUint64()
	if err != nil {
		log.Errorf("failed to decode buffer replacedId, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	var cmd Command
	cmdType := flag >> 6
	switch cmdType {
	case types.RequestType:
	case types.RequestOneWayType:
		cmd, ok = cachedCmd.(Request)
		if !ok {
			log.Errorf("cached cmd should be Request, maybe a bug occurred, contextId: %v", contextID)
			return types.StatusInternalFailure
		}

	case types.ResponseType:
		cmd, ok = cachedCmd.(Response)
		if !ok {
			log.Errorf("cached cmd should be Response, maybe a bug occurred, contextId: %v", contextID)
			return types.StatusInternalFailure
		}
	default:
		log.Errorf("failed to decode buffer, type = %s, value = %d", types.UnKnownRpcFlagType, flag)
		return types.StatusInternalFailure
	}

	id, err := buffer.ReadUint64()
	// we check encoded id equals cached command id
	if id != cmd.CommandId() {
		log.Errorf("encode buffer command id is not match , expect = %d, actual = %d", cmd.CommandId(), id)
		return types.StatusInternalFailure
	}

	// skip timeout or status
	buffer.ReadInt()

	dataBytes, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to decode buffer drain length, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	if dataBytes > 0 {
		cmd.SetData(Allocate(data[buffer.Pos():]))
	}

	// override cached request
	injectHeaderIfRequired(cmd, headers)

	// update command replacedId
	cmd.SetCommandId(replacedId)
	// call user extension implementation
	encode, err := ctx.Codec().Encode(context.TODO(), cmd)
	if err != nil {
		log.Fatalf("failed to encode command by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	attr.set(types.AttributeKeyEncodedBuffer, encode)

	// we don't need encode header again, the host side only pays attention to
	// the buffer of encode and sends it directly to the remote host
	proxyBuffer := encodeCommandBuffer(cmd, encode)
	// report encode data
	err = setEncodeBuffer(proxyBuffer.Bytes())
	if err != nil {
		log.Errorf("failed to report encode buffer by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
	}

	return types.StatusOK
}

//export proxy_keepalive_buffer_bytes
func proxyKeepAliveBufferBytes(contextID uint32, id uint64) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode keepalive buffer by protocol %s, contextId %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	cmd := ctx.KeepAlive().KeepAlive(id)
	attr := ctx.(attribute)
	attr.set(types.AttributeKeyEncodeCommand, cmd)

	return types.StatusOK
}

//export proxy_reply_keepalive_buffer_bytes
func proxyReplyKeepAliveBufferBytes(contextID uint32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode reply keepalive buffer by protocol %s, contextId %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	cmd := ctx.(attribute).attr(types.AttributeKeyDecodeCommand)

	//// todo how to obtain request ???
	//// convert data into an array of bytes to be parsed
	//data := parseByteSlice(bufferData, len)
	//buffer := Allocate(data)
	//cmd, err := ctx.Codec().Decode(context.TODO(), buffer)
	//if err != nil {
	//	log.Errorf("failed to decode reply keepalive request by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
	//	return types.StatusInternalFailure
	//}

	resp := ctx.KeepAlive().ReplyKeepAlive(cmd.(Request))
	attr := ctx.(attribute)
	attr.set(types.AttributeKeyEncodeCommand, resp)

	return types.StatusOK
}

//export proxy_hijack_buffer_bytes
func proxyHijackBufferBytes(contextID uint32, statusCode uint32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode hijack buffer by protocol %s, contextId %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	cmd := ctx.(attribute).attr(types.AttributeKeyDecodeCommand)

	//// todo how to obtain request ???
	//data := parseByteSlice(bufferData, len)
	//buffer := Allocate(data)
	//cmd, err := ctx.Codec().Decode(context.TODO(), buffer)
	//if err != nil {
	//	log.Errorf("failed to decode hijack request by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
	//	return types.StatusInternalFailure
	//}

	resp := ctx.Hijacker().Hijack(cmd.(Request), statusCode)
	attr := ctx.(attribute)
	attr.set(types.AttributeKeyEncodeCommand, resp)

	return types.StatusOK
}

func decodeCommandBuffer(cmd Command, drainBytes int) Buffer {
	// bufferData format:
	// encoded header map | Flag | Id | (Timeout|GetStatus) | drain length | raw bytes length | raw bytes
	headers := cmd.GetHeader()
	buf := AllocateBuffer()

	headerBytes := GetEncodeHeaderLength(headers)
	buf.WriteInt(headerBytes)
	// encoded header map
	if headerBytes > 0 {
		EncodeHeader(buf, headers)
	}

	var flag byte
	if cmd.IsHeartbeat() {
		flag = HeartBeatFlag
	}

	// should copy raw bytes
	flag = flag | CopyRawBytesFlag

	flagIndex := buf.Pos()
	// write flag
	buf.WriteByte(flag)
	// write id
	buf.WriteUint64(cmd.CommandId())

	// check is request
	if req, ok := cmd.(Request); ok {
		flag = flag | RpcRequestFlag
		if req.IsOneWay() {
			flag = flag | RpcOnewayFlag
		}
		// update request flag
		buf.PutByte(flagIndex, flag)
		buf.WriteUint32(req.GetTimeout())
	} else if resp, ok := cmd.(Response); ok {
		buf.WriteUint32(resp.GetStatus())
	}

	buf.WriteInt(drainBytes)
	if drainBytes > 0 {
		// write decode content length
		buf.WriteInt(cmd.GetData().Len())
		// write decode content, protocol header is not included
		buf.Write(cmd.GetData().Bytes())
	}

	return buf
}

func encodeCommandBuffer(cmd Command, encode Buffer) Buffer {
	// bufferData format:
	// encoded header map | Flag | Id | (Timeout|GetStatus) | drain length | raw bytes
	buf := AllocateBuffer()

	var headerBytes = 0
	buf.WriteInt(headerBytes)

	var flag byte
	if cmd.IsHeartbeat() {
		flag = HeartBeatFlag
	}

	// should copy raw bytes
	flag = flag | CopyRawBytesFlag

	flagIndex := buf.Pos()
	// write flag
	buf.WriteByte(flag)
	// write id
	buf.WriteUint64(cmd.CommandId())

	// check is request
	if req, ok := cmd.(Request); ok {
		flag = flag | RpcRequestFlag
		if req.IsOneWay() {
			flag = flag | RpcOnewayFlag
		}
		// update request flag
		buf.PutByte(flagIndex, flag)
		buf.WriteUint32(req.GetTimeout())
	} else if resp, ok := cmd.(Response); ok {
		buf.WriteUint32(resp.GetStatus())
	}

	drainBytes := encode.Len()
	buf.WriteInt(drainBytes)
	if drainBytes > 0 {
		buf.Write(encode.Bytes())
	}

	return buf
}

func injectHeaderIfRequired(cmd Command, headers *CommonHeader) {
	if cmd.GetHeader().Size() != headers.Size() {
		cmd.GetHeader().Range(func(key, value string) bool {
			v, ok := headers.Get(key)
			if !ok {
				// remove old key
				cmd.GetHeader().Del(key)
			} else {
				// add new key
				cmd.GetHeader().Set(key, v)
			}
			return true
		})
	}
}
