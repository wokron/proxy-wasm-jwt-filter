package config_test

import (
	"testing"
	"wasm-jwt-filter/pkg/config"
)

func TestParseConfig(t *testing.T) {
	data := []byte(`
	{
		"providers": [
			{
				"name": "123",
				"from_header": {
					"name": "Auth",
					"value_prefix": "XXX"
				},
				"from_param": "1234",
				"validate": {
					"issuer": "abc",
					"audiences": ["a", "b"],
					"jwk": "key"
				}
			},
			{
				"name": "1234",
				"validate": {
					"jwk": "key"
				}
			}
		],
		"rules": [
			{
				"match": {
					"path": "1234",
					"prefix": "123"
				},
				"requires": {
					"requires_any": [
						{
							"name": "abc"
						},
						{
							"requires_all": [
								{
									"name": "abc"
								},
								{
									"name": "abc"
								}
							]
						}
					]
				}
			},
			{
				"match": {
					
				}
			}
		]
	}
	`)

	config, err := config.ParseConfiguration(data)
	if err != nil {
		t.Error(err)
	}

	providers := config.Providers
	if len(providers) != 2 {
		t.Error(nil)
	}

	rules := config.Rules
	if len(rules) != 2 {
		t.Error(nil)
	}

	provider1 := providers[0]
	if provider1.Name != "123" || provider1.FromHeader == nil || provider1.FromParam == nil {
		t.Error(nil)
	}

	fromHeader1 := provider1.FromHeader
	if fromHeader1.Name != "Auth" || *fromHeader1.ValuePrefix != "XXX" {
		t.Error(nil)
	}

	validate1 := provider1.Validator
	if *validate1.Issuer != "abc" || validate1.Audiences == nil || len(*(validate1.Audiences)) != 2 || validate1.JWK != "key" {
		t.Error(nil)
	}

	provider2 := providers[1]
	if provider2.Name != "1234" || provider2.FromHeader != nil || provider2.FromParam != nil {
		t.Error(nil)
	}

	validate2 := provider2.Validator
	if validate2.Issuer != nil || validate2.Audiences != nil || validate2.JWK != "key" {
		t.Error(nil)
	}

	rule1 := rules[0]
	if rule1.Requires == nil {
		t.Error(nil)
	}

	match1 := rule1.Match
	if match1.Path == nil || *match1.Path != "1234" || match1.Prefix == nil || *match1.Prefix != "123" {
		t.Error(nil)
	}

	requires1 := rule1.Requires
	if requires1.RequiresAll != nil || requires1.NameRequired != nil || requires1.RequiresAny == nil {
		t.Error(nil)
	}

	requires1Any := *requires1.RequiresAny
	if *requires1Any[0].NameRequired != "abc" {
		t.Error(nil)
	}

	requires1AnyAll := requires1Any[1].RequiresAll
	if len(*requires1AnyAll) != 2 {
		t.Error(nil)
	}

	rule2 := rules[1]
	if rule2.Match.Path != nil || rule2.Match.Prefix != nil || rule2.Requires != nil {
		t.Error(nil)
	}
}
