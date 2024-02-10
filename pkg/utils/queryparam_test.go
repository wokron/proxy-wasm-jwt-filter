package utils_test

import (
	"testing"
	"wasm-jwt-filter/pkg/utils"
)

func TestGetQuryParamsFromPath(t *testing.T) {
	path := "/api/v1?a=1&b=2&c=&d&=true"
	params := utils.GetQueryParamsFromPath(path)

	if len(params) != 5 {
		t.Error(nil)
	}

	aVal, ok := params["a"]
	if !ok || aVal != "1" {
		t.Error(nil)
	}

	bVal, ok := params["b"]
	if !ok || bVal != "2" {
		t.Error(nil)
	}

	cVal, ok := params["c"]
	if !ok || cVal != "" {
		t.Error(nil)
	}

	dVal, ok := params["d"]
	if !ok || dVal != "" {
		t.Error(nil)
	}

	val, ok := params[""]
	if !ok || val != "true" {
		t.Error(nil)
	}
}
