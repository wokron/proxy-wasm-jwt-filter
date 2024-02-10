package main

import (
	"wasm-jwt-filter/pkg/config"
	"wasm-jwt-filter/pkg/jwt"
	"wasm-jwt-filter/pkg/rule"
	"wasm-jwt-filter/pkg/utils"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
	config *config.JWTFilterConfig
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("cannot get configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	config, err := config.ParseConfiguration(data)
	if err != nil {
		proxywasm.LogCriticalf("cannot parse configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	ctx.config = config

	counterIncrease := getCounterIncreaseFunc()
	jwt.ValidateCounter.IncreaseFunc = counterIncrease
	rule.RequestCounter.IncreaseFunc = counterIncrease

	return types.OnPluginStartStatusOK
}

func getCounterIncreaseFunc() func(label string, offset uint64) {
	counters := map[string]proxywasm.MetricCounter{}
	return func(label string, offset uint64) {
		counter, ok := counters[label]
		if !ok {
			counter = proxywasm.DefineCounterMetric(label)
			counters[label] = counter
			proxywasm.LogInfof("new metric label %s", label)
		}
		counter.Increment(offset)
	}
}

func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		config: ctx.config,
	}
}

type httpContext struct {
	types.DefaultHttpContext
	config *config.JWTFilterConfig
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	path, err := proxywasm.GetHttpRequestHeader(":path")
	if err != nil {
		proxywasm.LogErrorf("fail to get request header \":path\": %v", err)
		sendResponse(500, "Internal Server Error")
		panic(err)
	}

	accessHeader := func(key string) (string, bool) {
		value, err := proxywasm.GetHttpRequestHeader(key)
		return value, err == nil
	}
	accessQueryParam := func(key string) (string, bool) { return utils.GetQueryParamValue(path, key) }

	validateFunc := ctx.config.ProvidersMap.GetValidateFunc(path, accessHeader, accessQueryParam)
	ok, err := ctx.config.Rules.Validate(path, validateFunc)

	if err != nil {
		proxywasm.LogErrorf("plugin error: %v", err)
		sendResponse(500, "Internal Server Error\n")
		panic(err)
	}

	if !ok {
		proxywasm.LogInfo("forbidden: fail to validate all requirements")
		sendResponse(403, "Forbidden\n")
		return types.ActionPause
	}

	return types.ActionContinue
}

func sendResponse(statusCode uint32, body string) {
	err := proxywasm.SendHttpResponse(statusCode, nil, []byte(body), -1)
	if err != nil {
		proxywasm.LogErrorf("fail to send response, %v", err)
		panic(err)
	}
}
