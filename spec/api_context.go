package spec

import "github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"

//export proxy_on_context_create
func proxyOnContextCreate(contextID uint32, parentContextID uint32, contextType types.ContextType) uint32 {
	return 0
}

//export proxy_on_context_finalize
func proxyOnContextFinalize(contextID uint32) bool {
	return true
}
