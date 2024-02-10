package jwt

import (
	"strconv"
	"strings"

	"wasm-jwt-filter/pkg/metrics"
	"wasm-jwt-filter/pkg/rule"
	"wasm-jwt-filter/pkg/types"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang-jwt/jwt/v5"
)

var ValidateCounter = metrics.NewCounter("envoy_wasm_jwt_filter_validate_counts")

type ProvidersMap map[string]*JWTProvider

type MapAccessFunc func(string) (string, bool)

func (providersMap *ProvidersMap) GetValidateFunc(path string, accessHeader MapAccessFunc, accessQueryParam MapAccessFunc) rule.ValidateFunc {
	return func(name string) (bool, error) {
		provider, ok := (*providersMap)[name]
		if !ok {
			return false, types.ErrorInvalidConfig
		}

		return provider.Validate(path, accessHeader, accessQueryParam), nil
	}
}

type JWTProvider struct {
	Name       string       `json:"name"`
	FromHeader *JWTHeader   `json:"from_header"`
	FromParam  *string      `json:"from_param"`
	Validator  JWTValidator `json:"validate"`
}

func (provider *JWTProvider) Validate(path string, accessHeader MapAccessFunc, accessQueryParam MapAccessFunc) bool {
	jwtString, err := provider.ExtractJWTString(accessHeader, accessQueryParam)
	if err != nil {
		ValidateCounter.AddTag("provider_name", provider.Name).AddTag("success", strconv.FormatBool(false)).Increase(1)
		return false
	}

	ok := provider.Validator.Validate(jwtString)
	ValidateCounter.AddTag("provider_name", provider.Name).AddTag("success", strconv.FormatBool(ok)).Increase(1)
	return ok
}

func (provider *JWTProvider) ExtractJWTString(accessHeader MapAccessFunc, accessQueryParam MapAccessFunc) (string, error) {
	if provider.FromParam != nil {
		jwt, err := getJWTFromMap(accessQueryParam, *provider.FromParam, "")
		return jwt, err
	} else {
		key, prefix := provider.FromHeader.GetKeyAndPrefix()
		jwt, err := getJWTFromMap(accessHeader, key, prefix)
		return jwt, err
	}
}

func getJWTFromMap(accessMap MapAccessFunc, key string, prefix string) (string, error) {
	value, ok := accessMap(key)
	if !ok || !strings.HasPrefix(value, prefix) {
		return "", types.ErrorIllegalArgument
	}

	jwt := value[len(prefix):]
	return jwt, nil
}

type JWTHeader struct {
	Name        string  `json:"name"`
	ValuePrefix *string `json:"value_prefix"`
}

func (header *JWTHeader) GetKeyAndPrefix() (key string, prefix string) {
	if header == nil {
		key, prefix = "Authorization", "Bearer "
	} else {
		key = header.Name
		if header.ValuePrefix == nil {
			prefix = ""
		} else {
			prefix = *header.ValuePrefix
		}
	}
	return
}

type JWTValidator struct {
	Issuer    *string   `json:"issuer"`
	Audiences *[]string `json:"audiences"`
	JWK       string    `json:"jwk"`
}

func (requirement *JWTValidator) Validate(jwtString string) bool {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(jwtString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(requirement.JWK), nil
	})
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	ok := requirement.validateClaims(&claims)
	return ok
}

func (requirement *JWTValidator) validateClaims(claims *jwt.RegisteredClaims) bool {
	if requirement.Issuer != nil && *requirement.Issuer != claims.Issuer {
		return false
	}

	if requirement.Audiences != nil {
		expectAudiences := mapset.NewSet[string]((*requirement.Audiences)...)
		if !expectAudiences.ContainsAny([]string(claims.Audience)...) {
			return false
		}
	}

	return true
}
