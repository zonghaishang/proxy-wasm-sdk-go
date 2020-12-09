package types

// Action tell the host what action should be triggered
type Action uint32

const (
	Continue         Action = 1
	EndStream        Action = 2
	Done             Action = 3
	Pause            Action = 4
	WaitForMoreData  Action = 5
	WaitForEndOrFull Action = 6
	Close            Action = 7
)

// Status
type Status uint32

const (
	StatusOK                     Status = 0
	StatusEmpty                  Status = 1
	StatusNotFound               Status = 2
	StatusNotAllowed             Status = 3
	StatusBadArgument            Status = 4
	StatusInvalidMemoryAccess    Status = 5
	StatusInvalidOperation       Status = 6
	StatusCompareAndSwapMismatch Status = 7
)

type StreamType uint32

const (
	Downstream   StreamType = 1
	Upstream     StreamType = 2
	HttpRequest  StreamType = 3
	HttpResponse StreamType = 4
)

type ContextType uint32

const (
	VmContext     ContextType = 1
	PluginContext ContextType = 2
	StreamContext ContextType = 3
	HttpContext   ContextType = 4
)

type BufferType uint32

const (
	VmConfiguration         BufferType = 1
	PluginConfiguration     BufferType = 2
	DownstreamData          BufferType = 3
	UpstreamData            BufferType = 4
	HttpRequestBody         BufferType = 5
	HttpResponseBody        BufferType = 6
	HttpCalloutResponseBody BufferType = 7
)

type MapType uint32

const (
	HttpRequestHeaders       MapType = 1
	HttpRequestTrailers      MapType = 2
	HttpRequestMetadata      MapType = 3
	HttpResponseHeaders      MapType = 4
	HttpResponseTrailers     MapType = 5
	HttpResponseMetadata     MapType = 6
	HttpCallResponseHeaders  MapType = 7
	HttpCallResponseTrailers MapType = 8
	HttpCallResponseMetadata MapType = 9
	RpcRequestHeaders        MapType = 31
	RpcRequestTrailers       MapType = 32
	RpcResponseHeaders       MapType = 33
	RpcResponseTrailers      MapType = 34
)

// PeerType
type PeerType uint32

const (
	Local  PeerType = 1
	Remote PeerType = 2
)

// LogLevel proxy log level
type LogLevel uint32

const (
	LogLevelTrace LogLevel = 1
	LogLevelDebug LogLevel = 2
	LogLevelInfo  LogLevel = 3
	LogLevelWarn  LogLevel = 4
	LogLevelError LogLevel = 5
	LogLevelFatal LogLevel = 6
)

const (
	trace = "trace"
	debug = "debug"
	info  = "info"
	warn  = "warn"
	error = "error"
	fatal = "fatal"
)

func (level LogLevel) String() string {
	switch level {
	case LogLevelTrace:
		return trace
	case LogLevelDebug:
		return debug
	case LogLevelInfo:
		return info
	case LogLevelWarn:
		return warn
	case LogLevelError:
		return error
	case LogLevelFatal:
		return fatal
	default:
		panic("unsupported log level")
	}
}

type ExtensionType int

const (
	VmContextFilter     ExtensionType = 1
	PluginContextFilter ExtensionType = 2
	StreamContextFilter ExtensionType = 3
	HttpContextFilter   ExtensionType = 4
)
