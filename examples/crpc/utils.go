package crpc

import (
	"context"

	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
	"hash/fnv"
	"strings"
)

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func getUUID(bytes []byte) string {
	uid, err := FromBytes(bytes)
	if err != nil {
		proxy.Log.Warnf("uuid error %v", err)
		return string(bytes)
	}
	return uid.String()
}

func GetGovernValue(context context.Context, headers proxy.Header, key string) (string, bool) {
	if headers != nil {
		if val, ok := headers.Get(key); ok {
			return val, ok
		}
		if val, ok := headers.Get(strings.ToLower(key)); ok {
			return val, ok
		}
	}
	if context != nil {
		if val := context.Value(key); val != nil {
			return val.(string), true
		}
		if val := context.Value(strings.ToLower(key)); val != nil {
			return val.(string), true
		}
	}
	return "", false
}
