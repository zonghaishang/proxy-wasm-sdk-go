package proxy

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
)

//export proxy_decode_buffer_bytes
func proxyDecodeBufferBytes(contextID uint32, bufferData **byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode buffer by protocol %s, context id %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
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

	ctx.(Attribute).Set(types.AttributeKeyDecodeCommand, cmd)

	// we check decode is ok
	if err != nil {
		log.Fatalf("failed to decode buffer by protocol %s, context id %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	proxyBuffer := encodeProxyCommand(cmd, false)

	// report encode data
	err = setDecodeBuffer(proxyBuffer.Bytes())
	if err != nil {
		log.Errorf("failed to report decode buffer by protocol %s, context id %v, err %v", ctx.Name(), contextID, err)
	}

	return types.StatusOK
}

//export proxy_encode_buffer_bytes
func proxyEncodeBufferBytes(contextID uint32, bufferData **byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode buffer by protocol %s, context id %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}
	this.setActiveContextID(contextID)

	if len <= 0 {
		// should never be happen
		return types.StatusEmpty
	}

	// convert data into an array of bytes to be parsed
	data := parseByteSlice(*bufferData, len)
	buffer := Allocate(data)

	// bufferData format:
	// encoded header map | Flag | Id | (Timeout|Status) | drain length | raw bytes
	n, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to read decode buffer header map, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	headers := &CommonHeader{}
	// encoded header map
	if n > 0 {
		decodeHeader(data[4:4+n], headers)
	}

	flag, err := buffer.ReadByte()
	if err != nil {
		log.Errorf("failed to decode buffer flag, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	// find context cmd
	cachedCmd := ctx.(Attribute).Attr(types.AttributeKeyDecodeCommand)
	if cachedCmd == nil {
		// is heartbeat 、keep-alive or hijack ?
		cachedCmd = ctx.(Attribute).Attr(types.AttributeKeyEncodeCommand)
	}

	if cachedCmd == nil {
		log.Errorf("failed to find cached command, maybe a bug occurred, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	id, err := buffer.ReadUint64()
	if err != nil {
		log.Errorf("failed to decode buffer id, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	var cmd Command
	f := flag >> 6
	switch f {
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

	// todo check id
	//// we check encoded id equals cached command id
	//if id != cmd.CommandId() {
	//	log.Errorf("encode buffer command id is not match , expect = %d, actual = %d", cmd.CommandId(), id)
	//	return types.StatusInternalFailure
	//}

	// skip timeout or status
	buffer.ReadInt()

	bytes, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to decode buffer drain length, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	if bytes > 0 {
		// including protocol header
		cmd.SetData(Allocate(data[buffer.Pos():]))
	}

	// override cached request
	if cmd.Header().Size() != headers.Size() {
		cmd.Header().Range(func(key, value string) bool {
			v, ok := headers.Get(key)
			if !ok {
				// remove old key
				cmd.Header().Del(key)
			} else {
				// add new key
				cmd.Header().Set(key, v)
			}
			return true
		})
	}

	// update command id
	cmd.SetCommandId(id)
	// call user extension implementation
	buff, err := ctx.Codec().Encode(cmd)
	if err != nil {
		log.Fatalf("failed to encode command by protocol %s, context id %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}
	// update encoded command data
	cmd.SetData(buff)

	// we don't need encode header again, the host side only pays attention to
	// the buffer of encode and sends it directly to the remote host
	proxyBuffer := encodeProxyCommand(cmd, true)
	// report encode data
	err = setEncodeBuffer(proxyBuffer.Bytes())
	if err != nil {
		log.Errorf("failed to report encode buffer by protocol %s, context id %v, err %v", ctx.Name(), contextID, err)
	}

	return types.StatusOK
}

//export proxy_keepalive_buffer_bytes
func proxyKeepAliveBufferBytes(contextID uint32, id uint64) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode keepalive buffer by protocol %s, context id %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	cmd := ctx.KeepAlive().KeepAlive(id)
	attr := ctx.(Attribute)
	attr.Set(types.AttributeKeyEncodeCommand, cmd)

	return types.StatusOK
}

//export proxy_reply_keepalive_buffer_bytes
func proxyReplyKeepAliveBufferBytes(contextID uint32, bufferData **byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode reply keepalive buffer by protocol %s, context id %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	// convert data into an array of bytes to be parsed
	data := parseByteSlice(*bufferData, len)
	buffer := Allocate(data)
	cmd, err := ctx.Codec().Decode(buffer)
	if err != nil {
		log.Errorf("failed to decode reply keepalive request by protocol %s, context id %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	resp := ctx.KeepAlive().ReplyKeepAlive(cmd.(Request))
	attr := ctx.(Attribute)
	attr.Set(types.AttributeKeyEncodeCommand, resp)

	return types.StatusOK
}

//export proxy_hijack_buffer_bytes
func proxyHijackBufferBytes(contextID uint32, statusCode uint32, bufferData **byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode hijack buffer by protocol %s, context id %v not found", ctx.Name(), contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	data := parseByteSlice(*bufferData, len)
	buffer := Allocate(data)
	cmd, err := ctx.Codec().Decode(buffer)
	if err != nil {
		log.Errorf("failed to decode hijack request by protocol %s, context id %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	resp := ctx.Hijacker().Hijack(cmd.(Request), statusCode)
	attr := ctx.(Attribute)
	attr.Set(types.AttributeKeyEncodeCommand, resp)

	return types.StatusOK
}

func encodeProxyCommand(cmd Command, simple bool) Buffer {
	// bufferData format:
	// encoded header map | Flag | Id | (Timeout|Status) | drain length | raw bytes
	headers := cmd.Header()

	buf := AllocateBuffer()

	var n = 0
	if simple {
		buf.WriteInt(n)
	} else {
		n = getEncodeHeaderLength(headers)
		buf.WriteInt(n)
		// encoded header map
		if n > 0 {
			encodeHeader(buf, headers)
		}
	}

	var flag byte
	if cmd.IsHeartbeat() {
		flag = HeartBeatFlag
	}

	flagIndex := buf.Pos()
	// write flag
	buf.WriteByte(flag)
	// write id
	buf.WriteUint64(cmd.CommandId())

	// check is request
	if req, ok := cmd.(Request); ok {
		// update request flag
		buf.PutByte(flagIndex, flag|RpcRequestFlag)
		if req.IsOneWay() {
			buf.PutByte(flagIndex, flag|RpcOnewayFlag)
		}
		buf.WriteUint32(req.Timeout())
	} else if resp, ok := cmd.(Response); ok {
		buf.WriteUint32(resp.Status())
	}

	// update drain length
	buf.WriteInt(cmd.Data().Len())
	buf.Write(cmd.Data().Bytes())

	return buf
}
