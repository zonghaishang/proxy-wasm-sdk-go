package main

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/examples/crpc"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
)

func main() {
	proxy.SetNewRootContext(rootContext)
	proxy.SetNewProtocolContext(crpcContext)
}

type crpcProtocolContext struct {
	proxy.DefaultRootContext
	proxy.DefaultProtocolContext
	crpc      proxy.Protocol
	contextId uint32
}

func rootContext(rootContextId uint32) proxy.RootContext {
	return &crpcProtocolContext{
		crpc:      crpcProtocol,
		contextId: rootContextId,
	}
}

func crpcContext(rootContextId, contextId uint32) proxy.ProtocolContext {
	return &crpcProtocolContext{
		crpc:      crpcProtocol,
		contextId: rootContextId,
	}
}

var crpcProtocol = crpc.NewCrpcProtocol()

func (proto *crpcProtocolContext) Name() string {
	return proto.crpc.Name()
}
func (proto *crpcProtocolContext) Codec() proxy.Codec {
	return proto.crpc.Codec()
}

func (proto *crpcProtocolContext) KeepAlive() proxy.KeepAlive {
	return proto.crpc
}

func (proto *crpcProtocolContext) Hijacker() proxy.Hijacker {
	return proto.crpc
}

// vm and plugin lifecycle

func (proto *crpcProtocolContext) OnVMStart(conf proxy.ConfigMap) bool {

	proxy.Log.Infof("proxy_on_vm_start from Go!, config %v", conf)

	return true
}

func (proto *crpcProtocolContext) OnPluginStart(conf proxy.ConfigMap) bool {

	proxy.Log.Infof("proxy_on_plugin_start from Go!")

	return true
}
