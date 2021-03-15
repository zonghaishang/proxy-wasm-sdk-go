package main

import (
	"context"
	"github.com/zonghaishang/proxy-wasm-sdk-go/examples/crpc"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestCrpc(t *testing.T) {

	vmConfig := proxy.NewConfigMap()
	vmConfig.Set("engine", "wasm")

	opt := proxy.NewEmulatorOption().
		WithNewProtocolContext(crpcContext).
		WithNewRootContext(rootContext).
		WithVMConfiguration(vmConfig)

	host := proxy.NewHostEmulator(opt)
	// release lock and reset emulator state
	defer host.Done()
	// invoke host start vm
	host.StartVM()
	// invoke host plugin
	host.StartPlugin()

	// 1. invoke downstream decode
	ctxId := host.NewProtocolContext()
	// crpc plugin decode will be invoked
	cmd, err := host.Decode(ctxId, proxy.WrapBuffer(decodedRequestBytes(strconv.FormatUint(host.CurrentStreamId(), 10))))
	if err != nil {
		t.Fatalf("failed to invoke host decode request buffer, err: %v", err)
	}

	if _, ok := cmd.(*crpc.Request); !ok {
		t.Fatalf("decode request type error, expect *crpc.Request, actual %v", reflect.TypeOf(cmd))
	}

	// 2. invoke upstream encode
	upstreamBuf, err := host.Encode(ctxId, cmd)
	if err != nil {
		t.Fatalf("failed to invoke host encode request buffer, err: %v", err)
	}

	// check upstream content with downstream request
	if len(decodedRequestBytes(strconv.FormatUint(host.CurrentStreamId(), 10))) != len(upstreamBuf.Bytes()) {
		t.Fatalf("failed to invoke host encode request buffer, err: %v", err)
	}

	// complete protocol pipeline
	host.CompleteProtocolContext(ctxId)
}

func decodedRequestBytes(id string) []byte {
	rpcHeader := proxy.NewHeader()
	request := crpc.NewRequest(id, rpcHeader, proxy.WrapBuffer([]byte("crpc body")))
	request.Set(SERVICE_NAME_KEY, "test:1.0@crpc")
	request.Set(SERVICE_VERSION_KEY, "1.0")
	request.Set(GROUP_ID_KEY, "default")
	request.Set(SERVICE_METHOD_NAME_KEY, "sayHi")
	request.Set(TARGET_APP_NAME_KEY, "test1")
	request.Set(TRACE_ID_KEY, "0000000000")
	request.Set(SPAN_ID_KEY, "0000000002")
	request.Set(PARENT_SPAN_ID_KEY, "0000000001")
	request.Set(SAMPLED_KEY, "false")
	request.Set(FLAGS_KEY, "aaaa")
	request.Set(TRAN_NUM, "COBP20181105174343253620000000      ")

	request.CallApp = "1234"
	request.SourceApp = "1234"
	request.TranNum = "COBP20181105174343253620000000      "
	request.ApplySysTime = strconv.FormatInt(time.Now().UnixNano(), 10)
	baseStr := "00000000000000000000000000"
	if len(request.ApplySysTime) < 26 {
		request.ApplySysTime = baseStr[:26-len(request.ApplySysTime)] + request.ApplySysTime
	}
	if len(request.ApplySysTime) > 26 {
		request.ApplySysTime = request.ApplySysTime[:26]
	}
	crpc.SetRequestHeaderValue(request)
	buf, err := crpc.NewCrpcProtocol().Codec().Encode(context.TODO(), request)
	if err != nil {
		panic("failed to encode crpc request, err: " + err.Error())
	}
	return buf.Bytes()
}

const (
	SERVICE_NAME_KEY        = "serviceName"
	SERVICE_VERSION_KEY     = "serviceVersion"
	GROUP_ID_KEY            = "groupId"
	SERVICE_METHOD_NAME_KEY = "methodName"
	TARGET_APP_NAME_KEY     = "destinationAppName"
	TRACE_ID_KEY            = "traceId"
	SPAN_ID_KEY             = "spanId"
	PARENT_SPAN_ID_KEY      = "parentSpanId"
	SAMPLED_KEY             = "sampled"
	FLAGS_KEY               = "flags"
	TRAN_NUM                = "tranNum"
)
