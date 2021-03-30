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

	// crpc plugin decode will be invoked
	rsp, err := host.Decode(ctxId, proxy.WrapBuffer(decodeResponseBytes(strconv.FormatUint(host.CurrentStreamId(), 10))))
	if err != nil {
		t.Fatalf("failed to invoke host decode response buffer, err: %v", err)
	}

	if _, ok := rsp.(*crpc.Response); !ok {
		t.Fatalf("decode request type error, expect *crpc.Response, actual %v", reflect.TypeOf(cmd))
	}

	// 2. invoke upstream encode
	rspBuffer, err := host.Encode(ctxId, rsp)
	if err != nil {
		t.Fatalf("failed to invoke host encode response buffer, err: %v", err)
	}

	// check upstream content with downstream request
	if len(decodeResponseBytes(strconv.FormatUint(host.CurrentStreamId(), 10))) != len(rspBuffer.Bytes()) {
		t.Fatalf("failed to invoke host encode response buffer, err: %v", err)
	}

	// complete protocol pipeline
	host.CompleteProtocolContext(ctxId)
}

func decodeResponseBytes(id string) []byte {
	parseUint, _ := strconv.ParseUint(id, 10, 64)
	response := crpc.NewResponse(parseUint, "CRPC000", nil, nil)
	response.Set(SERVICE_NAME_KEY, "test:1.0@crpc")
	response.Set(SERVICE_VERSION_KEY, "1.0")
	response.Set(GROUP_ID_KEY, "default")
	response.Set(SERVICE_METHOD_NAME_KEY, "sayHi")
	response.Set(TARGET_APP_NAME_KEY, "test1")
	response.Set(TRACE_ID_KEY, "0000000000")

	response.RpcRespCode = "CRPC000"
	response.TranNum = "COBP20181105174343253620000000      "
	response.AppRespCode = "AAAAAAA"
	response.Heartbeat = false
	response.Body = proxy.NewBuffer(100)
	response.Body.WriteString("this is a test crpc response!")
	encode, _ := crpc.NewCrpcProtocol().Codec().Encode(context.TODO(), response)
	return encode.Bytes()

}

func decodedRequestBytes(id string) []byte {
	rpcHeader := proxy.NewHeader()
	parseUint, _ := strconv.ParseUint(id, 10, 64)
	request := crpc.NewRequest(parseUint, rpcHeader, proxy.WrapBuffer([]byte("crpc body")))
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

	commandId := request.CommandId()
	proxy.Log.Info(strconv.FormatUint(commandId, 10))

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
