package crpc

import (
	"encoding/binary"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"sync"
)

var requestIdMapper sync.Map

type RequestHeader struct {
	MagicNum         []byte //2 byte
	RequestLen       uint32 //4 byte
	HeaderLen        uint16 //2 byte
	Version          byte   //1 byte
	HeaderProperties []byte //3 byte
	Heartbeat        bool
	OneWay           bool
	RequestByte      []byte //16 byte
	RequestId        string //16 byte
	Timeout          uint32 //4 byte
	SourceApp        string //4 byte
	TranNum          string //36 byte
	ApplySysTime     string //26 byte
	CallApp          string //4 byte
	TagCnt           byte   //2 byte
	proxy.CommonHeader
}

func (h *RequestHeader) Clone() proxy.Header {
	clone := &RequestHeader{}
	*clone = *h
	clone.CommonHeader = *h.CommonHeader.Clone()
	return clone
}

type Request struct {
	RequestHeader

	rawData []byte // raw data
	rawTags []byte // sub slice of raw data, tag bytes
	rawBody []byte // sub slice of raw data, body bytes

	Data proxy.Buffer // wrapper of raw data
	Body proxy.Buffer // wrapper of raw body

	ContentChanged bool // indicate that content changed
}

func (r *Request) IsOneWay() bool {
	return r.OneWay
}

func (r *Request) GetTimeout() uint32 {
	return r.Timeout
}

func (r *Request) IsHeartbeat() bool {
	return r.Heartbeat
}

func (r *Request) CommandId() uint64 {
	return r.GetRequestId()
}

func (r *Request) SetCommandId(id uint64) {
	r.SetRequestId(id)
}

func (r *Request) GetRequestId() uint64 {
	return binary.BigEndian.Uint64(r.RequestByte[:8])
}

func (r *Request) SetRequestId(id uint64) {
	binary.BigEndian.PutUint64(r.RequestByte[:8], id)
	r.RequestId = getUUID(r.RequestByte[:])
	r.Set(GOVERN_REQUEST_ID, r.RequestId)
}

func (r *Request) IsHeartbeatFrame() bool {
	return r.RequestHeader.Heartbeat
}

func (r *Request) GetHeader() proxy.Header {
	return r
}

func (r *Request) GetData() proxy.Buffer {
	return r.Body
}

func (r *Request) SetData(data proxy.Buffer) {
	// judge if the address unchanged, assume that proxy logic will not operate the original Content buffer.
	if r.Body != data {
		r.ContentChanged = true
		r.Body = data
	}
}

type ResponseHeader struct {
	MagicNum         []byte //2 byte
	ResponseLen      uint32 //4 byte
	HeaderLen        uint16 //2 byte
	Version          byte   //1 byte
	HeaderProperties []byte //3 byte
	Heartbeat        bool
	RequestByte      []byte //16 byte
	RequestId        string //16 byte
	TranNum          string //36 byte
	RpcRespCode      string //7 byte
	AppRespCode      string //7 byte
	TagCnt           byte   //1 byte
	proxy.CommonHeader
}

func (h *ResponseHeader) Clone() proxy.Header {
	clone := &ResponseHeader{}
	*clone = *h

	// deep copy
	clone.CommonHeader = *h.CommonHeader.Clone()

	return clone
}

type Response struct {
	ResponseHeader
	rawData []byte
	rawTags []byte // sub slice of raw data, tag bytes
	rawBody []byte // sub slice of raw data, body bytes

	Data proxy.Buffer // wrapper of raw tags
	Body proxy.Buffer // wrapper of raw body

	ContentChanged bool
}

func (r *Response) GetStatus() uint32 {
	return r.GetStatusCode()
}

func (r *Response) IsHeartbeat() bool {
	return r.Heartbeat
}

func (r *Response) CommandId() uint64 {
	return r.GetRequestId()
}

func (r *Response) SetCommandId(id uint64) {
	r.SetRequestId(id)
}

// ~ XRespFrame
func (r *Response) GetRequestId() uint64 {
	return binary.BigEndian.Uint64(r.RequestByte[:8])
}

func (r *Response) SetRequestId(id uint64) {

	binary.BigEndian.PutUint64(r.RequestByte[:8], id)
	r.RequestId = getUUID(r.RequestByte[:])
	r.Set(GOVERN_REQUEST_ID, r.RequestId)
}

func (r *Response) IsHeartbeatFrame() bool {
	return r.ResponseHeader.Heartbeat
}

func (r *Response) GetHeader() proxy.Header {
	return r
}

func (r *Response) GetData() proxy.Buffer {
	return r.Body
}

func (r *Response) SetData(data proxy.Buffer) {
	// judge if the address unchanged, assume that proxy logic will not operate the original Content buffer.
	if r.Body != data {
		r.ContentChanged = true
		r.Body = data
	}
}

func (r *Response) GetStatusCode() uint32 {
	return 200
}
