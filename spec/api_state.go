package spec

import "github.com/zonghaishang/proxy-wasm-sdk-go/proxy"

type (
	HttpCalloutCallBack = func(headers proxy.Header, body proxy.Buffer)

	rootContextState struct {
		context       proxy.RootContext
		httpCallbacks map[uint32]*struct {
			callback        HttpCalloutCallBack
			callerContextID uint32
		}
	}
)

type state struct {
	newRootContext   func(contextID uint32) proxy.RootContext
	rootContexts     map[uint32]*rootContextState
	newFilterContext func(rootContextID, contextID uint32) proxy.FilterContext
	filterStreams    map[uint32]proxy.FilterContext
	newStreamContext func(rootContextID, contextID uint32) proxy.StreamContext
	streams          map[uint32]proxy.StreamContext
	newHttpContext   func(rootContextID, contextID uint32) proxy.HttpContext
	httpStreams      map[uint32]proxy.HttpContext

	// protocol context

	contextIDToRooID map[uint32]uint32
	activeContextID  uint32
}

var this = &state{
	rootContexts:     make(map[uint32]*rootContextState),
	httpStreams:      make(map[uint32]proxy.HttpContext),
	filterStreams:    make(map[uint32]proxy.FilterContext),
	streams:          make(map[uint32]proxy.StreamContext),
	contextIDToRooID: make(map[uint32]uint32),
}

func SetNewRootContext(f func(contextID uint32) proxy.RootContext) {
	this.newRootContext = f
}

func SetNewHttpContext(f func(rootContextID, contextID uint32) proxy.HttpContext) {
	this.newHttpContext = f
}

func SetNewFilterContext(f func(rootContextID, contextID uint32) proxy.FilterContext) {
	this.newFilterContext = f
}

func SetNewStreamContext(f func(rootContextID, contextID uint32) proxy.StreamContext) {
	this.newStreamContext = f
}

//go:inline
func (s *state) createRootContext(contextID uint32) {
	var ctx proxy.RootContext
	if s.newRootContext == nil {
		ctx = &proxy.DefaultRootContext{}
	} else {
		ctx = s.newRootContext(contextID)
	}

	s.rootContexts[contextID] = &rootContextState{
		context: ctx,
		httpCallbacks: map[uint32]*struct {
			callback        HttpCalloutCallBack
			callerContextID uint32
		}{},
	}
}

func (s *state) createFilterContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.filterStreams[contextID]; ok {
		panic("context id duplicated")
	}

	ctx := s.newFilterContext(rootContextID, contextID)
	s.contextIDToRooID[contextID] = rootContextID
	s.filterStreams[contextID] = ctx
}

func (s *state) createHttpContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.httpStreams[contextID]; ok {
		panic("context id duplicated")
	}

	ctx := s.newHttpContext(rootContextID, contextID)
	s.contextIDToRooID[contextID] = rootContextID
	s.httpStreams[contextID] = ctx
}

func (s *state) createStreamContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.streams[contextID]; ok {
		panic("context id duplicated")
	}

	ctx := s.newStreamContext(rootContextID, contextID)
	s.contextIDToRooID[contextID] = rootContextID
	s.streams[contextID] = ctx
}

func (s *state) registerHttpCallOut(calloutID uint32, callback HttpCalloutCallBack) {
	r := s.rootContexts[s.contextIDToRooID[s.activeContextID]]
	r.httpCallbacks[calloutID] = &struct {
		callback        HttpCalloutCallBack
		callerContextID uint32
	}{callback: callback, callerContextID: s.activeContextID}
}

//go:inline
func (s *state) setActiveContextID(contextID uint32) {
	s.activeContextID = contextID
}
