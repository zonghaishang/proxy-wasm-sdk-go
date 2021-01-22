package proxy

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
	"reflect"
)

type RootContext interface {
	OnVMStart(conf ConfigMap) bool
	OnPluginStart(conf ConfigMap) bool
	// OnTick()
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
	//// Context Used to save and pass data during a session
	//Context() context.Context
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

type Attribute interface {
	Attr(key string) interface{}
	Set(key string, v interface{})
	Remove(key string)
}

type (
	DefaultRootContext     struct{}
	DefaultFilterContext   struct{ DefaultAttribute }
	DefaultStreamContext   struct{}
	DefaultProtocolContext struct{ DefaultAttribute }
	DefaultAttribute       struct{ m map[string]interface{} }
)

var (
	_ RootContext     = &DefaultRootContext{}
	_ FilterContext   = &DefaultFilterContext{}
	_ StreamContext   = &DefaultStreamContext{}
	_ ProtocolContext = &DefaultProtocolContext{}
	_ Attribute       = &DefaultAttribute{}
)

// impl RootContext
//func (*DefaultRootContext) OnTick()                           {}
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

// attribute impl
func (a *DefaultAttribute) Attr(key string) interface{} {
	if a.m == nil {
		return nil
	}
	return a.m[key]
}

func (a *DefaultAttribute) Set(key string, v interface{}) {
	remove := v == nil || reflect.ValueOf(v).IsNil()

	if len(a.m) == 0 {
		// nil value should be ignored
		if remove {
			return
		}
		a.m = make(map[string]interface{})
	}

	if remove {
		a.Remove(key)
		return
	}

	a.m[key] = v
}

func (a *DefaultAttribute) Remove(key string) {
	if len(a.m) == 0 {
		return
	}
	// remove unused key
	delete(a.m, key)
}
