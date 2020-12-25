package proxy

type Codec interface {
	Decode(data Buffer) (Command, error)
	Encode(message Command) (Buffer, error)
}

// Command base request or response command
type Command interface {
	// Header get the data exchange header, maybe return nil.
	Header() Header
	// GetData return the complete message byte buffer, including the protocol header
	Data() Buffer
	// SetData update the complete message byte buffer, including the protocol header
	SetData(data Buffer)
	// IsHeartbeat check if the request is a heartbeat request
	IsHeartbeat() bool
	// CommandId get command id
	CommandId() uint64
	// SetCommandId update command id
	// In upstream, because of connection multiplexing,
	// the id of downstream needs to be replaced with id of upstream
	// blog: https://mosn.io/blog/posts/multi-protocol-deep-dive/#%E5%8D%8F%E8%AE%AE%E6%89%A9%E5%B1%95%E6%A1%86%E6%9E%B6
	SetCommandId(id uint64)
}

type Request interface {
	Command
	// IsOneWay Check that the request does not care about the response
	IsOneWay() bool
	Timeout() uint32 // request timeout
}

type Response interface {
	Command
	Status() uint32 // response status
}

type KeepAlive interface {
	KeepAlive(requestId uint64) Request
	ReplyKeepAlive(request Request) Response
}

type Hijacker interface {
	// Hijack allows sidecar to hijack requests
	Hijack(request Request, code uint32) Response
}
