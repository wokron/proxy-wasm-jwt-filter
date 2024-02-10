package config

import (
	"encoding/json"
	"wasm-jwt-filter/pkg/jwt"
	"wasm-jwt-filter/pkg/rule"
)

type JWTFilterConfig struct {
	Providers    []jwt.JWTProvider `json:"providers"`
	Rules        rule.RuleList     `json:"rules"`
	ProvidersMap jwt.ProvidersMap
}

func ParseConfiguration(data []byte) (*JWTFilterConfig, error) {
	config := JWTFilterConfig{ProvidersMap: map[string]*jwt.JWTProvider{}}
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	for _, provider := range config.Providers {
		config.ProvidersMap[provider.Name] = &provider
	}

	return &config, nil
}
