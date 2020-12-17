package proxy

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"
)

type RootContext interface {
	OnVMStart(conf ConfigMap) bool
	OnVMDone() bool
	OnPluginStart(conf ConfigMap) bool
	OnTick()
}

type StreamContext interface {
	OnDownStreamReceived(headers Header, buffer Buffer, trailers Header)
	OnUpstreamReceived(headers Header, buffer Buffer, trailers Header)
	OnDownstreamClose(peerType types.PeerType)
	OnNewConnection() types.Action
	OnUpstreamClose(peerType types.PeerType)
	OnStreamDone()
}

type HttpContext interface {
	OnHttpRequestReceived(headers Header, body Buffer)
	OnHttpResponseReceived(headers Header, body Buffer)
	OnHttpStreamDone()
}
