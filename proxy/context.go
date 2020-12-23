package proxy

import (
	"context"
	"github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"
)

type RootContext interface {
	OnVMStart(conf ConfigMap) bool
	OnPluginStart(conf ConfigMap) bool
	OnTick()
	OnVMDone() bool
}

// L7 layer extension
type FilterContext interface {
	// OnDownStreamReceived Called when the data requests,
	// The caller should check if the parameter value is nil
	OnDownStreamReceived(headers Header, buffer Buffer, trailers Header) types.Action
	// OnUpstreamReceived Called when the data responds,
	// The caller should check if the parameter value is nil
	OnUpstreamReceived(headers Header, buffer Buffer, trailers Header) types.Action
	// Context Used to save and pass data during a session
	Context() context.Context
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

type ProtocolContext interface {
}

type (
	DefaultRootContext   struct{}
	DefaultFilterContext struct {
		ctx context.Context
	}
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

// impl FilterContext
func (c *DefaultFilterContext) OnDownStreamReceived(headers Header, buffer Buffer, trailers Header) types.Action {
	return types.ActionContinue
}

func (c *DefaultFilterContext) OnUpstreamReceived(headers Header, buffer Buffer, trailers Header) types.Action {
	return types.ActionContinue
}

func (c *DefaultFilterContext) Context() context.Context {
	if c.ctx == nil {
		c.ctx = &internalContext{Context: context.Background()}
	}
	return c.ctx
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

// context impl
type ContextKey int

const (
	ContextKeyStreamID ContextKey = iota
	ContextKeyListenerType
	ContextKeyHeaderHolder
	ContextKeyBufferHolder
	ContextKeyTrailerHolder
	ContextKeyEnd
)

type internalContext struct {
	context.Context
	values [ContextKeyEnd]interface{}
}

func (c *internalContext) Value(key interface{}) interface{} {
	if contextKey, ok := key.(ContextKey); ok {
		return c.values[contextKey]
	}
	return c.Context.Value(key)
}

func Get(ctx context.Context, key ContextKey) interface{} {
	if context, ok := ctx.(*internalContext); ok {
		return context.values[key]
	}
	return ctx.Value(key)
}

func WithValue(parent context.Context, key ContextKey, value interface{}) context.Context {
	if context, ok := parent.(*internalContext); ok {
		context.values[key] = value
		return context
	}
	context := &internalContext{Context: parent}
	context.values[key] = value
	return context
}
