package main

import (
	"context"
	"github.com/zonghaishang/proxy-wasm-sdk-go/examples/bolt"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"reflect"
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

	// invoke downstream decode
	ctxId := host.NewProtocolContext()
	// bolt plugin decode will be invoked
	cmd, err := host.Decode(ctxId, proxy.WrapBuffer(decodedRequestBytes()))
	if err != nil {
		t.Fatalf("failed to invoke host decode request buffer, err: %v", err)
	}

	if _, ok := cmd.(*bolt.Request); !ok {
		t.Fatalf("decode request type error, expect *bolt.Request, actual %v", reflect.TypeOf(cmd))
	}

}

func decodedRequestBytes() []byte {
	rpcHeader := proxy.NewHeader()
	rpcHeader.Set("service", "com.alipay.demo.HelloService")
	request := bolt.NewRpcRequest(1, rpcHeader, proxy.WrapBuffer([]byte("bolt body")))
	buf, err := bolt.NewBoltProtocol().Codec().Encode(context.TODO(), request)
	if err != nil {
		panic("failed to encode bolt request, err: " + err.Error())
	}
	return buf.Bytes()
}
