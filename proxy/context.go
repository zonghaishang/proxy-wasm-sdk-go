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
	Options() Options     // protocol options
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
	DefaultRootContext     struct{}
	DefaultFilterContext   struct{ ctx context.Context }
	DefaultStreamContext   struct{}
	DefaultProtocolContext struct{}
)

var (
	_ RootContext     = &DefaultRootContext{}
	_ FilterContext   = &DefaultFilterContext{}
	_ StreamContext   = &DefaultStreamContext{}
	_ ProtocolContext = &DefaultProtocolContext{}
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

// impl ProtocolContext
func (*DefaultProtocolContext) Name() string {
	panic("protocol name should be override")
}

func (*DefaultProtocolContext) Codec() Codec {
	panic("protocol codec should be override")
}

func (*DefaultProtocolContext) KeepAlive() KeepAlive {
	panic("protocol keepalive should be override")
}

func (*DefaultProtocolContext) Hijacker() Hijacker {
	panic("protocol hijacker should be override")
}

func (*DefaultProtocolContext) Options() Options {
	return options
}
