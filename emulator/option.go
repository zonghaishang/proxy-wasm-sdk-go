package emulator

import "github.com/zonghaishang/proxy-wasm-sdk-go/proxy"

type Option struct {
	pluginConfiguration proxy.ConfigMap
	vmConfiguration     proxy.ConfigMap
	newRootContext      func(uint32) proxy.RootContext
	newStreamContext    func(uint32, uint32) proxy.StreamContext
	newFilterContext    func(uint32, uint32) proxy.FilterContext
	newProtocolContext  func(uint32, uint32) proxy.ProtocolContext
}

func NewEmulatorOption() *Option {
	return &Option{}
}

func (o *Option) WithNewRootContext(f func(uint32) proxy.RootContext) *Option {
	o.newRootContext = f
	return o
}

func (o *Option) WithNewHttpContext(f func(uint32, uint32) proxy.FilterContext) *Option {
	o.newFilterContext = f
	return o
}

func (o *Option) WithNewStreamContext(f func(uint32, uint32) proxy.StreamContext) *Option {
	o.newStreamContext = f
	return o
}

func (o *Option) WithNewProtocolContext(f func(uint32, uint32) proxy.ProtocolContext) *Option {
	o.newProtocolContext = f
	return o
}

func (o *Option) WithPluginConfiguration(data proxy.ConfigMap) *Option {
	o.pluginConfiguration = data
	return o
}

func (o *Option) WithVMConfiguration(data proxy.ConfigMap) *Option {
	o.vmConfiguration = data
	return o
}
