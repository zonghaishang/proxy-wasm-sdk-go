package crpc

import "time"

const (
	PROTOCOL_NAME           = "crpc" //中信银行私有 rpc 协议
	ProtocolFirstByte  byte = 0x1A
	ProtocolSecondByte byte = 0x19

	DefaultExpiration      = 15 * time.Second
	DefaultCleanupInterval = 60 * time.Second

	HeaderBeginLen int = 6

	RequestHeaderLen          int = 102 // protocol header fields length
	RequestHeartBeatHeaderLen int = 28  // heartbeat header length
	RequestTagStartIndex      int = 103
	RequestHeaderLenBeforeTag     = RequestTagStartIndex - HeaderBeginLen - 2

	ResponseTagStartIndex      int = 79
	ResponseHeaderLenBeforeTag     = ResponseTagStartIndex - HeaderBeginLen - 2

	HeartBeatRequestHeaderLen                     int = 22
	HeartBeatHeaderLen                                = HeartBeatRequestHeaderLen - 2
	SOFARPC_ROUTER_RULE_METADATA_CRPC_TARGET_HOST     = "crpc_target_host"
	SOFARPC_ROUTER_RULE_METADATA_CRPC_TARGET_ZONE     = "crpc_target_zone"
	EGRESS                                            = "egress"
	INGRESS                                           = "ingress"

	VarProxyTryTimeout    string = "proxy_try_timeout"
	VarProxyGlobalTimeout string = "proxy_global_timeout"
	VarProxyHijackStatus  string = "proxy_hijack_status"
	VarProxyGzipSwitch    string = "proxy_gzip_switch"

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

	CRPC_SUCCESS            = "CRPC000"
	CRPC_RPC_REQUEST_ERROR  = "CRPC001"
	CRPC_RPC_RESPONSE_ERROR = "CRPC002"
	CRPC_TIMEOUT            = "CRPC003"
	CRPC_ROUTER_ERROR       = "CRPC004"
	CRPC_ERROR              = "CRPC005"

	MRPC_LIMIT_ERROR      = "MRPC001"
	MRPC_CIRCUIT_ERROR    = "MRPC002"
	MRPC_AUTH_ERROR       = "MRPC003"
	MRPC_DOWNGROUD_ERROR  = "MRPC004"
	MRPC_ROUTE_ERROR      = "MRPC005"
	MRPC_FAULT_RULE_ERROR = "MRPC006"
	MRPC_UNKNOWN_ERROR    = "MRPC009"
	//非服务治理相关的响应码
	MRPC_REQUEST_ERROR_CRPC001  = "MRPC101"
	MRPC_RESPONSE_ERROR_CRPC002 = "MRPC102"
	MRPC_TIMEOUT_CRPC003        = "MRPC103"

	AuthHiJackCode          = 403
	RouteHiJackCode         = 404
	CircuitBreakHiJackCode  = 418
	LimitExceededHiJackCode = 429
	FaultInjectHiJackCode   = 520
	DownGradeHiJackCode     = 509

	StatusOK                  = 200
	StatusGatewayTimeout      = 504
	StatusBadRequest          = 400
	StatusInternalServerError = 500

	GOVERN_TEST_TRAFFIC_KEY     = "X-Govern-Test-Traffic"
	GOVERN_SERVICE_KEY          = "X-Govern-Service"
	GOVERN_SERVICE_TYPE_KEY     = "X-Govern-Service-Type"
	GOVERN_METHOD_KEY           = "X-Govern-Method"
	GOVERN_TIMEOUT_KEY          = "X-Govern-Timeout"
	GOVERN_SOURCE_APP_KEY       = "X-Govern-Source-App"
	GOVERN_SOURCE_IP_KEY        = "X-Govern-Source-Ip"
	GOVERN_TARGET_APP_KEY       = "X-Govern-Target-App"
	GOVERN_RESP_CODE_KEY        = "X-Govern-Resp-Code"
	GOVERN_HIJACK_CODE_KEY      = "X-Govern-HiJack-Code"
	GOVERN_REQUEST_ID           = "X-Request-Id"
	GOVERN_TRACE_ID             = "X-B3-TraceId"
	GOVERN_TRACE_SPAN_ID        = "X-B3-SpanId"
	GOVERN_TRACE_PARENT_SPAN_ID = "X-B3-ParentSpanId"
	GOVERN_TRACE_SAMPLED        = "X-B3-Sampled"
	GOVERN_TRACE_flags          = "X-B3-Flags"

	CRPC_TRACER_HEADER_TARGET_DATACENTER  = "crpc_target_datacenter"
	CRPC_TRACER_HEADER_SOURCE_ZONE        = "crpc_source_zone"
	CRPC_TRACER_HEADER_SOURCE_DATACENTER  = "crpc_source_datacenter"
	CRPC_TRACER_HEADER_RESPONSE_CRPC_CODE = "response_rpc_code"
	CRPC_TRACER_HEADER_RESPONSE_BIZ_CODE  = "response_biz_code"
)
