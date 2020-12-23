package proxy

import (
	"fmt"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy/types"
)

var log = SetLogger(NewLogger())

type proxyLogger struct {
}

func NewLogger() Logger {
	return &proxyLogger{}
}

func (p *proxyLogger) Debug(msg string) {
	ABI_ProxyLog(types.LogLevelDebug, stringBytePtr(msg), len(msg))
}

func (p *proxyLogger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ABI_ProxyLog(types.LogLevelDebug, stringBytePtr(msg), len(msg))
}

func (p *proxyLogger) Info(msg string) {
	ABI_ProxyLog(types.LogLevelInfo, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ABI_ProxyLog(types.LogLevelInfo, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Warn(msg string) {
	ABI_ProxyLog(types.LogLevelWarn, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ABI_ProxyLog(types.LogLevelWarn, stringBytePtr(msg), len(msg))
}

func (p *proxyLogger) Error(msg string) {
	ABI_ProxyLog(types.LogLevelError, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ABI_ProxyLog(types.LogLevelError, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Fatal(msg string) {
	ABI_ProxyLog(types.LogLevelFatal, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ABI_ProxyLog(types.LogLevelFatal, stringBytePtr(msg), len(msg))
}

//export proxy_on_log
func proxyOnLog(contextID uint32) {
	if ctx, ok := this.streams[contextID]; ok {
		this.setActiveContextID(contextID)
		ctx.OnLog()
	} else if ctx, ok := this.filterStreams[contextID]; ok {
		this.setActiveContextID(contextID)
		ctx.OnLog()
	} else if ctx, ok := this.rootContexts[contextID]; ok {
		this.setActiveContextID(contextID)
		ctx.context.OnLog()
	} else {
		panic("invalid context on proxy_on_done")
	}
}
