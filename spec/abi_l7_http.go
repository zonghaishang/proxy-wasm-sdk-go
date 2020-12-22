package spec

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"
)

//export proxy_on_request_headers
func proxyOnRequestHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	ctx, ok := this.httpStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_headers")
	}
	this.setActiveContextID(contextID)
	if endOfStream && numHeaders > 0 {
		hs, err := GetHttpRequestHeaders()
		if err != nil {
			log.Errorf("failed to get request headers: %v", err)
			return types.ActionContinue
		}
		return ctx.OnHttpRequestReceived(proxy.CommonHeader(hs), proxy.NewBuffer(0))
	}
	return types.ActionContinue
}

//export proxy_on_request_body
func proxyOnRequestBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	ctx, ok := this.httpStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_body")
	}
	this.setActiveContextID(contextID)

	if endOfStream && bodySize > 0 {
		hs, err := GetHttpRequestHeaders()
		if err != nil {
			log.Errorf("failed to get request headers: %v", err)
			return types.ActionContinue
		}

		body, err := GetHttpRequestBody(0, bodySize)
		if err != nil {
			log.Errorf("failed to get request body: %v", err)
			return types.ActionContinue
		}
		return ctx.OnHttpRequestReceived(proxy.CommonHeader(hs), proxy.Allocate(body))
	}
	return types.ActionContinue
}

//export proxy_on_response_headers
func proxyOnResponseHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	ctx, ok := this.httpStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)

	if endOfStream && numHeaders > 0 {
		hs, err := GetHttpResponseHeaders()
		if err != nil {
			log.Errorf("failed to get request headers: %v", err)
			return types.ActionContinue
		}
		return ctx.OnHttpResponseReceived(proxy.CommonHeader(hs), proxy.NewBuffer(0))
	}
	return types.ActionContinue
}

//export proxy_on_response_body
func proxyOnResponseBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	ctx, ok := this.httpStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)

	if endOfStream && bodySize > 0 {
		hs, err := GetHttpResponseHeaders()
		if err != nil {
			log.Errorf("failed to get response headers: %v", err)
			return types.ActionContinue
		}

		body, err := GetHttpResponseBody(0, bodySize)
		if err != nil {
			log.Errorf("failed to get response body: %v", err)
			return types.ActionContinue
		}
		return ctx.OnHttpRequestReceived(proxy.CommonHeader(hs), proxy.Allocate(body))
	}
	return types.ActionContinue
}
