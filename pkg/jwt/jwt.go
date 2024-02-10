package jwt

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang-jwt/jwt/v5"
	"wasm-jwt-filter/pkg/types"
)

type JWTProvider struct {
	Name       string       `json:"name"`
	FromHeader *JWTHeader   `json:"from_header"`
	FromParam  *string      `json:"from_param"`
	Validator  JWTValidator `json:"validate"`
}

type MapAccessFunc func(string) (string, error)

func (provider *JWTProvider) ExtractJWTString(accessHeader MapAccessFunc, accessQueryParam MapAccessFunc) (string, error) {
	if provider.FromHeader != nil {
		key, prefix := provider.FromHeader.GetKeyAndPrefix()
		jwt, err := getJWTFromMap(accessHeader, key, prefix)
		return jwt, err
	} else if provider.FromParam != nil {
		jwt, err := getJWTFromMap(accessQueryParam, *provider.FromParam, "")
		return jwt, err
	} else {
		return "", types.ErrorInvalidConfig
	}
}

func getJWTFromMap(accessMap MapAccessFunc, key string, prefix string) (string, error) {
	value, err := accessMap(key)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(value, prefix) {
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

func (requirement *JWTValidator) ValidateRequirement(jwtString string) bool {
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
