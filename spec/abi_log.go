package spec

import (
	"fmt"
	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"github.com/zonghaishang/proxy-wasm-sdk-go/spec/types"
)

var log = proxy.SetLogger(NewLogger())

type proxyLogger struct {
}

func NewLogger() proxy.Logger {
	return &proxyLogger{}
}

func (p *proxyLogger) Debug(msg string) {
	ProxyLog(types.LogLevelDebug, stringBytePtr(msg), len(msg))
}

func (p *proxyLogger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ProxyLog(types.LogLevelDebug, stringBytePtr(msg), len(msg))
}

func (p *proxyLogger) Info(msg string) {
	ProxyLog(types.LogLevelInfo, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ProxyLog(types.LogLevelInfo, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Warn(msg string) {
	ProxyLog(types.LogLevelWarn, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ProxyLog(types.LogLevelWarn, stringBytePtr(msg), len(msg))
}

func (p *proxyLogger) Error(msg string) {
	ProxyLog(types.LogLevelError, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ProxyLog(types.LogLevelError, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Fatal(msg string) {
	ProxyLog(types.LogLevelFatal, stringBytePtr(msg), len(msg))
}
func (p *proxyLogger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ProxyLog(types.LogLevelFatal, stringBytePtr(msg), len(msg))
}
