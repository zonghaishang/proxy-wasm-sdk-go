package proxy

import (
	"fmt"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
	stdout "log"
)

type (
	protocolEmulator struct {
		protocolStreams map[uint32]*protocolStreamState
		streamId        uint64
	}
	protocolStreamState struct {
		decodedCmd      Command // decoded command
		decodedProxyBuf Buffer  // report to host buffer

		encodedBuf      Buffer // encoded buffer
		encodedProxyBuf Buffer // report to host buffer

		request        Request  // keepalive request
		response       Response // keepalive response
		hijackRequest  Request  // hijack request
		hijackResponse Response // hijack response
		Status         types.Status
	}
)

func newProtocolEmulator() *protocolEmulator {
	host := &protocolEmulator{protocolStreams: map[uint32]*protocolStreamState{}, streamId: 1}
	return host
}

// protocol L7 level
func (h *protocolEmulator) NewProtocolContext() (contextID uint32) {
	contextID = getNextContextID()
	proxyOnContextCreate(contextID, RootContextID)
	h.protocolStreams[contextID] = &protocolStreamState{Status: types.StatusOK}
	return
}

func (h *protocolEmulator) CurrentStreamId() uint64 {
	return h.streamId
}

func (h *protocolEmulator) Decode(contextID uint32, data Buffer) (Command, error) {
	cs, ok := h.protocolStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	bufferData := &data.Bytes()[0]
	cs.Status = proxyDecodeBufferBytes(contextID, bufferData, data.Len())

	if cs.Status == types.StatusOK {
		return cs.decodedCmd, nil
	}

	return nil, fmt.Errorf("decode error, code %d", cs.Status)
}

func (h *protocolEmulator) Encode(contextID uint32, cmd Command) (Buffer, error) {
	cs, ok := h.protocolStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	buf := AllocateBuffer()
	// encode data format:
	// encoded header map | Flag | replaceId, id | (Timeout|GetStatus) | drain length | raw dataBytes
	headerBytes := 0
	if cmd.GetHeader().Size() > 0 {
		headerBytes = GetEncodeHeaderLength(cmd.GetHeader())
	}
	// encode header map
	buf.WriteInt(headerBytes)
	// encoded header map
	if headerBytes > 0 {
		EncodeHeader(buf, cmd.GetHeader())
	}

	flagIndex := buf.Len()
	// should copy raw bytes
	flag := CopyRawBytesFlag
	if cmd.IsHeartbeat() {
		flag = flag | HeartBeatFlag
	}
	buf.WriteByte(flag)

	// generate stream id
	h.streamId += 2
	// write replaced id
	buf.WriteUint64(h.streamId)
	// write command id
	buf.WriteUint64(cmd.CommandId())

	if req, ok := cmd.(Request); ok {
		flag = flag | RpcRequestFlag
		if req.IsOneWay() {
			flag = flag | RpcOnewayFlag
		}
		// update request flag
		buf.PutByte(flagIndex, flag)
		// write timeout
		buf.WriteUint32(req.GetTimeout())
	} else if resp, ok := cmd.(Response); ok {
		// write status code
		buf.WriteUint32(resp.GetStatus())
	}

	dataBytes := cmd.GetData().Len()
	// write drain length
	buf.WriteInt(dataBytes)
	if dataBytes > 0 {
		// write raw dataBytes
		buf.Write(cmd.GetData().Bytes())
	}

	// invoke the plugin encode
	cs.Status = proxyEncodeBufferBytes(contextID, &buf.Bytes()[0], buf.Len())
	if cs.Status == types.StatusOK {
		return cs.encodedBuf, nil
	}

	return nil, fmt.Errorf("encode error, code %d", cs.Status)
}

// heartbeat
func (h *protocolEmulator) KeepAlive(contextID uint32, requestId uint64) Request {
	cs, ok := h.protocolStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	cs.Status = proxyKeepAliveBufferBytes(contextID, int64(requestId))
	if cs.Status == types.StatusOK {
		cs.request = h.protocolEmulatorProxyKeepAlive()
		return cs.request
	}

	return nil
}

func (h *protocolEmulator) ReplyKeepAlive(contextID uint32, request Request) Response {
	cs, ok := h.protocolStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	// todo: buffer data format should be
	// encoded header map | Flag | Id | (Timeout|GetStatus) | drain length | raw bytes
	bufferData := request.GetData().Bytes()[0]
	cs.Status = proxyReplyKeepAliveBufferBytes(contextID, &bufferData, request.GetData().Len())
	if cs.Status == types.StatusOK {
		cs.response = h.protocolEmulatorProxyReplyKeepAlive()
		return cs.response
	}

	return nil
}

// hijacker
func (h *protocolEmulator) Hijack(contextID uint32, request Request, code uint32) Response {
	cs, ok := h.protocolStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	// todo: buffer data format should be
	// encoded header map | Flag | Id | (Timeout|GetStatus) | drain length | raw bytes
	bufferData := request.GetData().Bytes()[0]
	cs.Status = proxyHijackBufferBytes(contextID, int32(code), &bufferData, request.GetData().Len())
	if cs.Status == types.StatusOK {
		cs.response = h.protocolEmulatorProxyReplyKeepAlive()
		return cs.response
	}

	return nil
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *protocolEmulator) protocolEmulatorProxySetBufferBytes(bt types.BufferType, start int, maxSize int,
	bufferData *byte, bufferSize int) types.Status {
	body := parseByteSlice(bufferData, bufferSize)
	active := VMStateGetActiveContextID()
	stream := h.protocolStreams[active]
	ctx := this.protocolStreams[active]
	switch bt {
	case types.BufferTypeDecodeData:
		stream.decodedProxyBuf = WrapBuffer(body)
		stream.decodedCmd = ctx.(attribute).attr(types.AttributeKeyDecodeCommand).(Command)
		// stream.decodedBuf =
	case types.BufferTypeEncodeData:
		stream.encodedProxyBuf = WrapBuffer(body)
		stream.encodedBuf = ctx.(attribute).attr(types.AttributeKeyEncodedBuffer).(Buffer)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
	return types.StatusOK
}

func (h *protocolEmulator) protocolEmulatorProxyKeepAlive() Request {
	active := VMStateGetActiveContextID()
	ctx := this.protocolStreams[active]
	return ctx.(attribute).attr(types.AttributeKeyEncodeCommand).(Request)
}

func (h *protocolEmulator) protocolEmulatorProxyReplyKeepAlive() Response {
	active := VMStateGetActiveContextID()
	ctx := this.protocolStreams[active]
	return ctx.(attribute).attr(types.AttributeKeyEncodeCommand).(Response)
}

// impl HostEmulator
func (h *protocolEmulator) CompleteProtocolContext(contextID uint32) {
	proxyOnLog(contextID)
	proxyOnDone(contextID)
	proxyOnDelete(contextID)
}
