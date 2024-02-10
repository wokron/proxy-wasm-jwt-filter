package main

import (
	"testing"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestJWTFilter(t *testing.T) {
	config := `
	{
		"providers": [
		  	{
				"name": "provider-1",
				"validate": {
			  		"jwk": "your-secure-key"
				}
		  	}
		],
		"rules": [
			{
				"match": {
					"path": "/api/v1/abc"
				}
			},
			{
				"match": {
					"prefix": "/api/v1"
				},
				"requires": {
					"name": "provider-1"
				}
			}
		]
	  }
	`

	vmTest(t, func(t *testing.T, v types.VMContext) {
		opt := proxytest.NewEmulatorOption().
			WithPluginConfiguration([]byte(config)).
			WithVMContext(v)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		if types.OnPluginStartStatusOK != host.StartPlugin() {
			t.Error(nil)
		}

		contextID := host.InitializeHttpContext()

		headers1 := [][2]string{{":path", "/api/v1/abc"}}
		acion := host.CallOnRequestHeaders(contextID, headers1, false)
		if acion != types.ActionContinue {
			t.Error(nil)
		}

		headers2 := [][2]string{{":path", "/api/v1/abcd"}}
		acion = host.CallOnRequestHeaders(contextID, headers2, false)
		if acion != types.ActionPause {
			t.Error(nil)
		}

		headers3 := [][2]string{{":path", "/api/v1/abcd"}, {"Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkZW1vIiwiaWF0IjoxNzA3NjI2NjUzLCJuYmYiOjE3MDc2MjY2NTN9.bm3fS837e3214XxeSWAcfH0ZuRRJsl2kKnU4tIqPfgg"}}
		acion = host.CallOnRequestHeaders(contextID, headers3, false)
		if acion != types.ActionContinue {
			t.Error(nil)
		}

		headers4 := [][2]string{{":path", "/api/v1/abcd"}, {"Authorization", "Bearer invalid.jwt"}}
		acion = host.CallOnRequestHeaders(contextID, headers4, false)
		if acion != types.ActionPause {
			t.Error(nil)
		}

		headers5 := [][2]string{{":path", "/api/v2/abcd"}, {"Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkZW1vIiwiaWF0IjoxNzA3NjI2NjUzLCJuYmYiOjE3MDc2MjY2NTN9.bm3fS837e3214XxeSWAcfH0ZuRRJsl2kKnU4tIqPfgg"}}
		acion = host.CallOnRequestHeaders(contextID, headers5, false)
		if acion != types.ActionPause {
			t.Error(nil)
		}

		value, err := host.GetCounterMetric("envoy_wasm_jwt_filter_request_counts_match_no=0_permit=true") // 1
		if err != nil || value != 1 {
			t.Error(nil)
		}

		value, err = host.GetCounterMetric("envoy_wasm_jwt_filter_request_counts_match_no=1_permit=true") // 3
		if err != nil || value != 1 {
			t.Error(nil)
		}

		value, err = host.GetCounterMetric("envoy_wasm_jwt_filter_request_counts_match_no=1_permit=false") // 2, 4
		if err != nil || value != 2 {
			t.Error(err)
		}

		value, err = host.GetCounterMetric("envoy_wasm_jwt_filter_request_counts_match_no=-1_permit=false") // 5
		if err != nil || value != 1 {
			t.Error(err)
		}

		value, err = host.GetCounterMetric("envoy_wasm_jwt_filter_validate_counts_provider_name=provider-1_success=true") // 3
		if err != nil || value != 1 {
			t.Error(err)
		}

		value, err = host.GetCounterMetric("envoy_wasm_jwt_filter_validate_counts_provider_name=provider-1_success=false") // 2, 4
		if err != nil || value != 2 {
			t.Error(err)
		}
	})

}

func vmTest(t *testing.T, f func(*testing.T, types.VMContext)) {
	t.Helper()

	t.Run("go", func(t *testing.T) {
		f(t, &vmContext{})
	})
}
