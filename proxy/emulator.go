package proxy

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
	"sync"
)

type HostEmulator interface {
	// release emulator resource and reset state
	Done()

	// Root
	StartVM()
	StartPlugin()
	FinishVM()

	GetLogs(level types.LogLevel) []string
	GetTickPeriod() uint32
	Tick()

	// protocol L7 level
	NewProtocolContext() (contextID uint32)
	Decode(contextID uint32, data Buffer) (Command, error)
	Encode(contextID uint32, cmd Command) (Buffer, error)
	// heartbeat
	KeepAlive(contextID uint32, requestId uint64) Request
	ReplyKeepAlive(contextID uint32, request Request) Response
	// hijacker
	Hijack(contextID uint32, request Request, code uint32) Response

	// filter L7 level
	NewFilterContext() (contextID uint32)
	PutRequestHeaders(contextID uint32, headers map[string]string, endOfStream bool)
	PutRequestBody(contextID uint32, body []byte, endOfStream bool)
	PutRequestTrailers(contextID uint32, headers map[string]string)

	GetRequestHeaders(contextID uint32) (headers map[string]string)
	GetRequestBody(contextID uint32) []byte

	PutResponseHeaders(contextID uint32, headers map[string]string, endOfStream bool)
	PutResponseBody(contextID uint32, body []byte, endOfStream bool)
	PutResponseTrailers(contextID uint32, headers map[string]string)

	GetResponseHeaders(contextID uint32) (headers map[string]string)
	GetResponseBody(contextID uint32) []byte

	CompleteHttpStream(contextID uint32)
	GetCurrentStreamAction(contextID uint32) types.Action

	// network
}

var (
	hostMux       = sync.Mutex{}
	nextContextID = RootContextID + 1
)

type hostEmulator struct {
	*rootEmulator
	*protocolEmulator
	*networkEmulator
	*filterEmulator

	effectiveContextID uint32
}
