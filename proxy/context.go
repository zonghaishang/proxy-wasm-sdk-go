package proxy

import (
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
	OnDownStreamReceived(headers Header, buffer Buffer, trailers Header) types.Action
	OnUpstreamReceived(headers Header, buffer Buffer, trailers Header) types.Action
}

// L7 layer extension
type HttpContext interface {
	OnHttpRequestReceived(headers Header, body Buffer) types.Action
	OnHttpResponseReceived(headers Header, body Buffer) types.Action
	OnHttpStreamDone()
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
	DefaultFilterContext struct{}
	DefaultHttpContext   struct{}
	DefaultStreamContext struct{}
)

var (
	_ RootContext   = &DefaultRootContext{}
	_ FilterContext = &DefaultFilterContext{}
	_ HttpContext   = &DefaultHttpContext{}
	_ StreamContext = &DefaultStreamContext{}
)

// impl RootContext
func (*DefaultRootContext) OnTick()                           {}
func (*DefaultRootContext) OnVMStart(conf ConfigMap) bool     { return true }
func (*DefaultRootContext) OnPluginStart(conf ConfigMap) bool { return true }
func (*DefaultRootContext) OnVMDone() bool                    { return true }

// impl FilterContext
func (*DefaultFilterContext) OnDownStreamReceived(headers Header, buffer Buffer, trailers Header) types.Action {
	return types.ActionContinue
}

func (*DefaultFilterContext) OnUpstreamReceived(headers Header, buffer Buffer, trailers Header) types.Action {
	return types.ActionContinue
}

// impl HttpContext
func (*DefaultHttpContext) OnHttpRequestReceived(headers Header, body Buffer) types.Action {
	return types.ActionContinue
}
func (*DefaultHttpContext) OnHttpResponseReceived(headers Header, body Buffer) types.Action {
	return types.ActionContinue
}
func (*DefaultHttpContext) OnHttpStreamDone() {}

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
