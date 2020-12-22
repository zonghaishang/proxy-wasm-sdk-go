package spec

import "github.com/zonghaishang/proxy-wasm-sdk-go/proxy"

//export proxy_on_vm_start
func proxyOnVMStart(rootContextID uint32, vmConfigurationSize int) bool {
	ctx, ok := this.rootContexts[rootContextID]
	if !ok {
		panic("invalid context on proxy_on_vm_start")
	}
	this.setActiveContextID(rootContextID)
	configBytes, err := GetVMConfiguration(vmConfigurationSize)
	if err != nil {
		log.Errorf("failed to get vm config, error: %s", err.Error())
		return false
	}

	return ctx.context.OnVMStart(proxy.CommonHeader(DecodeMap(configBytes)))
}

//export proxy_on_plugin_start
func proxyOnPluginStart(rootContextID uint32, pluginConfigurationSize int) bool {
	ctx, ok := this.rootContexts[rootContextID]
	if !ok {
		panic("invalid context on proxy_on_configure")
	}
	this.setActiveContextID(rootContextID)
	configBytes, err := GetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		log.Errorf("failed to get plugin config, error: %s", err.Error())
		return false
	}
	return ctx.context.OnPluginStart(proxy.CommonHeader(DecodeMap(configBytes)))
}
