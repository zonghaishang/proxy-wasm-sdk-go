package proxy

const (
	HeartBeatFlag  byte = 1 << 5      // 0010_0000
	RpcRequestFlag byte = 1 << 6      // 0100_0000
	RpcOnewayFlag  byte = 1<<6 | 1<<7 // 1100_0000
)

//export proxy_on_context_create
func proxyOnContextCreate(contextID uint32, rootContextID uint32) {
	if rootContextID == 0 {
		this.createRootContext(contextID)
	} else if this.newFilterContext != nil {
		this.createFilterContext(contextID, rootContextID)
	} else if this.newProtocolContext != nil {
		this.createProtocolContext(contextID, rootContextID)
	} else if this.newStreamContext != nil {
		this.createStreamContext(contextID, rootContextID)
	} else {
		panic("invalid context id on proxy_on_context_create")
	}
}

//export proxy_on_done
func proxyOnDone(contextID uint32) bool {
	if ctx, ok := this.filterStreams[contextID]; ok {
		this.setActiveContextID(contextID)
		ctx.OnFilterStreamDone()
		return true
	} else if ctx, ok := this.streams[contextID]; ok {
		this.setActiveContextID(contextID)
		ctx.OnStreamDone()
		return true
	} else if ctx, ok := this.rootContexts[contextID]; ok {
		this.setActiveContextID(contextID)
		response := ctx.context.OnVMDone()
		return response
	} else {
		panic("invalid context on proxy_on_done")
	}
}

//export proxy_on_delete
func proxyOnDelete(contextID uint32) {
	if _, ok := this.filterStreams[contextID]; ok {
		delete(this.filterStreams, contextID)
	} else if _, ok := this.streams[contextID]; ok {
		delete(this.streams, contextID)
	} else if _, ok = this.rootContexts[contextID]; ok {
		delete(this.rootContexts, contextID)
	} else {
		panic("invalid context on proxy_on_done")
	}
	delete(this.contextIDToRootID, contextID)
}

//type internalContext struct {
//	context.Context
//	values [types.ContextKeyEnd]interface{}
//}
//
//func (c *internalContext) Value(key interface{}) interface{} {
//	if contextKey, ok := key.(types.ContextKey); ok {
//		return c.values[contextKey]
//	}
//	return c.Context.Value(key)
//}
//
//func Get(ctx context.Context, key types.ContextKey) interface{} {
//	if context, ok := ctx.(*internalContext); ok {
//		return context.values[key]
//	}
//	return ctx.Value(key)
//}
//
//func WithValue(parent context.Context, key types.ContextKey, value interface{}) context.Context {
//	if context, ok := parent.(*internalContext); ok {
//		context.values[key] = value
//		return context
//	}
//	context := &internalContext{Context: parent}
//	context.values[key] = value
//	return context
//}
