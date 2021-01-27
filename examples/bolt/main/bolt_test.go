package main

import (
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"testing"
)

func TestBolt(t *testing.T) {

	vmConfig := proxy.NewConfigMap()
	vmConfig.Set("engine", "wasm")

	opt := proxy.NewEmulatorOption().
		WithNewProtocolContext(boltContext).
		WithNewRootContext(rootContext).
		WithVMConfiguration(vmConfig)

	host := proxy.NewHostEmulator(opt)
	// release lock and reset emulator state
	defer host.Done()
	// invoke host start vm
	host.StartVM()
	// invoke host plugin
	host.StartPlugin()

}
