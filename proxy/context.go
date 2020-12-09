package proxy

import "github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"

type RootContext interface {
	OnVMStart(conf map[string]string) bool
	OnVMDone() bool
	OnPluginStart(conf map[string]string) bool
	OnTick()
}

type StreamContext interface {
	OnDownStreamReceived(headers map[string]string, buf []byte, trailers map[string]string)
	OnUpstreamReceived(headers map[string]string, buf []byte, trailers map[string]string)
	OnDownstreamClose(peerType types.PeerType)
	OnNewConnection() types.Action
	OnUpstreamClose(peerType types.PeerType)
	OnStreamDone()
}

type HttpContext interface {
	OnHttpRequestReceived(headers map[string]string, body []byte, trailers map[string]string)
	OnHttpResponseReceived(headers map[string]string, body []byte, trailers map[string]string)
	OnHttpStreamDone()
}
