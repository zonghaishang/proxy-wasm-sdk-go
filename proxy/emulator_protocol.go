package proxy

import (
	"fmt"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
	stdout "log"
)

type (
	protocolEmulator struct {
		protocolStreams map[uint32]*protocolStreamState
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

// protocol L7 level
func (h *protocolEmulator) NewProtocolContext() (contextID uint32) {
	contextID = getNextContextID()
	proxyOnContextCreate(contextID, RootContextID)
	h.protocolStreams[contextID] = &protocolStreamState{Status: types.StatusOK}
	return
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

	bufferData := cmd.GetData().Bytes()
	cs.Status = proxyEncodeBufferBytes(contextID, &bufferData[0], cmd.GetData().Len())

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

	cs.Status = proxyKeepAliveBufferBytes(contextID, requestId)
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
	cs.Status = proxyHijackBufferBytes(contextID, code, &bufferData, request.GetData().Len())
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
		stream.decodedProxyBuf = Allocate(body)
		stream.decodedCmd = ctx.(attribute).attr(types.AttributeKeyDecodeCommand).(Command)
		// stream.decodedBuf =
	case types.BufferTypeEncodeData:
		stream.encodedProxyBuf = Allocate(body)
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
