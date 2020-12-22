package spec

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"
)

//export proxy_on_new_connection
func proxyOnNewConnection(contextID uint32) types.Action {
	ctx, ok := this.streams[contextID]
	if !ok {
		panic("invalid context")
	}
	this.setActiveContextID(contextID)
	return ctx.OnNewConnection()
}

//export proxy_on_downstream_data
func proxyOnDownstreamData(contextID uint32, dataSize int, endOfStream bool) types.Action {
	ctx, ok := this.streams[contextID]
	if !ok {
		panic("invalid context")
	}
	this.setActiveContextID(contextID)

	if dataSize == 0 {
		return types.ActionContinue
	}

	data, err := GetDownStreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		log.Errorf("failed to get downstream data: %v", err)
		return types.ActionContinue
	}

	return ctx.OnDownstreamData(proxy.Allocate(data), endOfStream)
}

//export proxy_on_downstream_connection_close
func proxyOnDownstreamConnectionClose(contextID uint32, pType types.PeerType) {
	ctx, ok := this.streams[contextID]
	if !ok {
		panic("invalid context")
	}
	this.setActiveContextID(contextID)
	ctx.OnDownstreamClose(pType)
}

//export proxy_on_upstream_data
func proxyOnUpstreamData(contextID uint32, dataSize int, endOfStream bool) types.Action {
	ctx, ok := this.streams[contextID]
	if !ok {
		panic("invalid context")
	}
	this.setActiveContextID(contextID)

	if dataSize == 0 {
		return types.ActionContinue
	}

	data, err := GetUpstreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		log.Errorf("failed to get upstream data: %v", err)
		return types.ActionContinue
	}

	return ctx.OnUpstreamData(proxy.Allocate(data), endOfStream)
}

//export proxy_on_upstream_connection_close
func proxyOnUpstreamConnectionClose(contextID uint32, pType types.PeerType) {
	ctx, ok := this.streams[contextID]
	if !ok {
		panic("invalid context")
	}
	this.setActiveContextID(contextID)
	ctx.OnUpstreamClose(pType)
}
