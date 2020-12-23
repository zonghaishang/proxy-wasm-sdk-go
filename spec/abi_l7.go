package spec

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"
)

//export proxy_on_request_headers
func proxyOnRequestHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_headers")
	}
	this.setActiveContextID(contextID)

	var header proxy.Header
	if numHeaders > 0 {
		hs, err := GetHttpRequestHeaders()
		if err != nil {
			log.Errorf("failed to get request headers: %v", err)
			return types.ActionContinue
		}
		header = proxy.CommonHeader(hs)
		// update context header
		proxy.WithValue(ctx.Context(), proxy.ContextKeyHeaderHolder, header)
	}

	if endOfStream {
		return ctx.OnDownStreamReceived(header, proxy.NewBuffer(0), nil)
	}

	return types.ActionContinue
}

//export proxy_on_request_body
func proxyOnRequestBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_body")
	}
	this.setActiveContextID(contextID)

	var body proxy.Buffer
	if bodySize > 0 {
		bodyBytes, err := GetHttpRequestBody(0, bodySize)
		if err != nil {
			log.Errorf("failed to get request body: %v", err)
			return types.ActionContinue
		}

		body = proxy.Allocate(bodyBytes)
		// update context body buffer
		proxy.WithValue(ctx.Context(), proxy.ContextKeyBufferHolder, body)
	}

	if endOfStream {
		header := proxy.Get(ctx.Context(), proxy.ContextKeyHeaderHolder)
		return ctx.OnDownStreamReceived(header.(proxy.Header), body, nil)
	}

	return types.ActionContinue
}

//export proxy_on_request_trailers
func proxyOnRequestTrailers(contextID uint32, numTrailers int) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_trailers")
	}
	this.setActiveContextID(contextID)

	var trailer proxy.Header
	if numTrailers > 0 {
		trailers, err := GetHttpRequestTrailers()
		if err != nil {
			log.Errorf("failed to get request trailer: %v", err)
			return types.ActionContinue
		}
		trailer = proxy.CommonHeader(trailers)
		// update context header
		proxy.WithValue(ctx.Context(), proxy.ContextKeyTrailerHolder, trailer)
	}

	header := proxy.Get(ctx.Context(), proxy.ContextKeyHeaderHolder)
	body := proxy.Get(ctx.Context(), proxy.ContextKeyBufferHolder)

	return ctx.OnDownStreamReceived(header.(proxy.Header), body.(proxy.Buffer), trailer)
}

//export proxy_on_response_headers
func proxyOnResponseHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)

	var header proxy.Header
	if numHeaders > 0 {
		hs, err := GetHttpResponseHeaders()
		if err != nil {
			log.Errorf("failed to get response headers: %v", err)
			return types.ActionContinue
		}
		header = proxy.CommonHeader(hs)
		// update context header
		proxy.WithValue(ctx.Context(), proxy.ContextKeyHeaderHolder, header)
	}

	if endOfStream {
		return ctx.OnUpstreamReceived(header, proxy.NewBuffer(0), nil)
	}

	return types.ActionContinue
}

//export proxy_on_response_body
func proxyOnResponseBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)

	var body proxy.Buffer
	if bodySize > 0 {
		bodyBytes, err := GetHttpResponseBody(0, bodySize)
		if err != nil {
			log.Errorf("failed to get response body: %v", err)
			return types.ActionContinue
		}

		body = proxy.Allocate(bodyBytes)
		// update context body buffer
		proxy.WithValue(ctx.Context(), proxy.ContextKeyBufferHolder, body)
	}

	if endOfStream {
		header := proxy.Get(ctx.Context(), proxy.ContextKeyHeaderHolder)
		return ctx.OnUpstreamReceived(header.(proxy.Header), body, nil)
	}

	return types.ActionContinue
}

//export proxy_on_response_trailers
func proxyOnResponseTrailers(contextID uint32, numTrailers int) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)
	var trailer proxy.Header
	if numTrailers > 0 {
		trailers, err := GetHttpResponseTrailers()
		if err != nil {
			log.Errorf("failed to get request trailer: %v", err)
			return types.ActionContinue
		}
		trailer = proxy.CommonHeader(trailers)
		// update context header
		proxy.WithValue(ctx.Context(), proxy.ContextKeyTrailerHolder, trailer)
	}

	header := proxy.Get(ctx.Context(), proxy.ContextKeyHeaderHolder)
	body := proxy.Get(ctx.Context(), proxy.ContextKeyBufferHolder)

	return ctx.OnUpstreamReceived(header.(proxy.Header), body.(proxy.Buffer), trailer)
}
