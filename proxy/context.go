package proxy

import (
	"context"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
)

type RootContext interface {
	OnVMStart(conf ConfigMap) bool
	OnPluginStart(conf ConfigMap) bool
	OnTick()
	OnVMDone() bool
	OnLog()
}

// L7 layer extension
type FilterContext interface {
	// OnDownStreamReceived Called when the data requests,
	// The caller should check if the parameter value is nil
	OnDownStreamReceived(headers Header, data Buffer, trailers Header) types.Action
	// OnUpstreamReceived Called when the data responds,
	// The caller should check if the parameter value is nil
	OnUpstreamReceived(headers Header, data Buffer, trailers Header) types.Action
	// Context Used to save and pass data during a session
	Context() context.Context
	OnFilterStreamDone()
	OnLog()
}

// L7 layer extension
type ProtocolContext interface {
	Name() string         // protocol name
	Codec() Codec         // frame encode & decode
	KeepAlive() KeepAlive // protocol keep alive
	Hijacker() Hijacker   // protocol hijacker
}

type Codec interface {
	Decode(data Buffer) (Command, error)
	Encode(message Command) (Buffer, error)
}

// Command base request or response command
type Command interface {
	// Header get the data exchange header, maybe return nil.
	Header() Header
	// GetData return the complete message byte buffer, including the protocol header
	Data() Buffer
	// SetData update the complete message byte buffer, including the protocol header
	SetData(data Buffer)
	// IsHeartbeat check if the request is a heartbeat request
	IsHeartbeat() bool
	// CommandId get command id
	CommandId() uint64
	// SetCommandId update command id
	// In upstream, because of connection multiplexing,
	// the id of downstream needs to be replaced with id of upstream
	// blog: https://mosn.io/blog/posts/multi-protocol-deep-dive/#%E5%8D%8F%E8%AE%AE%E6%89%A9%E5%B1%95%E6%A1%86%E6%9E%B6
	SetCommandId(id uint64)
}

type Request interface {
	Command
	// IsOneWay Check that the request does not care about the response
	IsOneWay() bool
	Timeout() uint32 // request timeout
}

type Response interface {
	Command
	Status() uint32 // response status
}

type KeepAlive interface {
	KeepAlive(requestId uint64) Request
	ReplyKeepAlive(request Request) Response
}

type Hijacker interface {
	// Hijack allows sidecar to hijack requests
	Hijack(request Request, code uint32) Response
}

// L4 layer extension (host not support now.)
type StreamContext interface {
	OnDownstreamData(buffer Buffer, endOfStream bool) types.Action
	OnDownstreamClose(peerType types.PeerType)
	OnNewConnection() types.Action
	OnUpstreamData(buffer Buffer, endOfStream bool) types.Action
	OnUpstreamClose(peerType types.PeerType)
	OnStreamDone()
	OnLog()
}

type (
	DefaultRootContext   struct{}
	DefaultFilterContext struct{ ctx context.Context }
	DefaultStreamContext struct{}
)

var (
	_ RootContext   = &DefaultRootContext{}
	_ FilterContext = &DefaultFilterContext{}
	_ StreamContext = &DefaultStreamContext{}
)

// impl RootContext
func (*DefaultRootContext) OnTick()                           {}
func (*DefaultRootContext) OnVMStart(conf ConfigMap) bool     { return true }
func (*DefaultRootContext) OnPluginStart(conf ConfigMap) bool { return true }
func (*DefaultRootContext) OnVMDone() bool                    { return true }
func (*DefaultRootContext) OnLog() {
}

// impl FilterContext
func (c *DefaultFilterContext) OnDownStreamReceived(headers Header, data Buffer, trailers Header) types.Action {
	return types.ActionContinue
}

func (c *DefaultFilterContext) OnUpstreamReceived(headers Header, data Buffer, trailers Header) types.Action {
	return types.ActionContinue
}

func (c *DefaultFilterContext) Context() context.Context {
	if c.ctx == nil {
		c.ctx = &internalContext{Context: context.Background()}
	}
	return c.ctx
}

func (*DefaultFilterContext) OnFilterStreamDone() {}

func (*DefaultFilterContext) OnLog() {
}

// impl StreamContext
func (*DefaultStreamContext) OnDownstreamData(buffer Buffer, endOfStream bool) types.Action {
	return types.ActionContinue
}

func (*DefaultStreamContext) OnDownstreamClose(peerType types.PeerType) {
}

func (*DefaultStreamContext) OnNewConnection() types.Action {
	return types.ActionContinue
}

func (*DefaultStreamContext) OnUpstreamData(buffer Buffer, endOfStream bool) types.Action {
	return types.ActionContinue
}

func (*DefaultStreamContext) OnUpstreamClose(peerType types.PeerType) {
}

func (*DefaultStreamContext) OnStreamDone() {
}

func (*DefaultStreamContext) OnLog() {
}
