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

type StreamContext interface {
	OnDownStreamReceived(headers Header, buffer Buffer, trailers Header) types.Action
	OnDownstreamClose(peerType types.PeerType)
	OnUpstreamReceived(headers Header, buffer Buffer, trailers Header) types.Action
	OnNewConnection() types.Action
	OnUpstreamClose(peerType types.PeerType)
	OnStreamDone()
}

type HttpContext interface {
	OnHttpRequestReceived(headers Header, body Buffer) types.Action
	OnHttpResponseReceived(headers Header, body Buffer) types.Action
	OnHttpStreamDone()
}

type ProtocolContext interface {
}

type (
	DefaultRootContext   struct{}
	DefaultStreamContext struct{}
	DefaultHttpContext   struct{}
)

var (
	_ RootContext   = &DefaultRootContext{}
	_ StreamContext = &DefaultStreamContext{}
	_ HttpContext   = &DefaultHttpContext{}
)

// impl RootContext
func (*DefaultRootContext) OnTick()                           {}
func (*DefaultRootContext) OnVMStart(conf ConfigMap) bool     { return true }
func (*DefaultRootContext) OnPluginStart(conf ConfigMap) bool { return true }
func (*DefaultRootContext) OnVMDone() bool                    { return true }

// impl StreamContext
func (*DefaultStreamContext) OnDownStreamReceived(headers Header, buffer Buffer, trailers Header) types.Action {
	return types.ActionContinue
}
func (*DefaultStreamContext) OnDownstreamClose(peerType types.PeerType) {}
func (*DefaultStreamContext) OnUpstreamReceived(headers Header, buffer Buffer, trailers Header) types.Action {
	return types.ActionContinue
}
func (*DefaultStreamContext) OnNewConnection() types.Action           { return types.ActionContinue }
func (*DefaultStreamContext) OnUpstreamClose(peerType types.PeerType) {}
func (*DefaultStreamContext) OnStreamDone()                           {}

// impl HttpContext
func (*DefaultHttpContext) OnHttpRequestReceived(headers Header, body Buffer) types.Action {
	return types.ActionContinue
}
func (*DefaultHttpContext) OnHttpResponseReceived(headers Header, body Buffer) types.Action {
	return types.ActionContinue
}
func (*DefaultHttpContext) OnHttpStreamDone() {}
