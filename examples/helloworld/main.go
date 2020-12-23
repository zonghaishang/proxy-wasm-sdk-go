// Copyright 2020 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"math/rand"
	"time"

	"github.com/zonghaishang/proxy-wasm-sdk-go/proxy"
)

const tickMilliseconds uint32 = 1000

func main() {
	proxy.SetNewRootContext(newHelloWorld)
}

type helloWorld struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxy.DefaultRootContext
	contextID uint32
}

func newHelloWorld(contextID uint32) proxy.RootContext {
	return &helloWorld{contextID: contextID}
}

// override
func (ctx *helloWorld) OnVMStart(conf proxy.ConfigMap) bool {
	rand.Seed(time.Now().UnixNano())

	proxy.Log.Info("proxy_on_vm_start from Go!")
	if err := proxy.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxy.Log.Fatalf("failed to set tick period: %v", err)
	}

	return true
}

// override
func (ctx *helloWorld) OnTick() {
	t := time.Now().UnixNano()
	proxy.Log.Infof("It's %d: random value: %d", t, rand.Uint64())
}
