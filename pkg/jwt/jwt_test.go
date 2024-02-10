package jwt_test

import (
	"testing"
	"wasm-jwt-filter/pkg/jwt"
)

func TestGetKeyAndPrefix(t *testing.T) {
	prefixStr := "XXX"
	header1 := jwt.JWTHeader{
		Name:        "Auth",
		ValuePrefix: &prefixStr,
	}

	key, prefix := header1.GetKeyAndPrefix()
	if key != "Auth" || prefix != "XXX" {
		t.Error(nil)
	}

	header2 := jwt.JWTHeader{
		Name: "Auth",
	}

	key, prefix = header2.GetKeyAndPrefix()
	if key != "Auth" || prefix != "" {
		t.Error(nil)
	}

	var header3 *jwt.JWTHeader = nil
	key, prefix = header3.GetKeyAndPrefix()
	if key != "Authorization" || prefix != "Bearer " {
		t.Error(nil)
	}
}

func TestJWTValidator(t *testing.T) {
	issuer := "abc"
	audiences := []string{"a", "b"}

	validator1 := jwt.JWTValidator{
		JWK: "key1",
	}

	jwt1 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDc1NTc3OTMsIm5iZiI6MTcwNzU1Nzc5M30.Zze-emcAhn5qsn5EjdGM5DHZ-Dq7jAmUHWQRe94m6eQ"
	ok := validator1.Validate(jwt1)
	if !ok {
		t.Error(nil)
	}

	jwt1 = "invalid jwt"
	ok = validator1.Validate(jwt1)
	if ok {
		t.Error(nil)
	}

	validator2 := jwt.JWTValidator{
		JWK:    "key1",
		Issuer: &issuer,
	}

	// has issuer
	jwt2 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDc1NTc3OTMsIm5iZiI6MTcwNzU1Nzc5MywiaXNzIjoiYWJjIn0.6rUwgjQhi0ThZVedAkPFEihrc4qSBJ1xWcMJkpb4Kqg"
	ok = validator2.Validate(jwt2)
	if !ok {
		t.Error(nil)
	}

	// has no issuer
	jwt2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDc1NTc3OTMsIm5iZiI6MTcwNzU1Nzc5M30.Zze-emcAhn5qsn5EjdGM5DHZ-Dq7jAmUHWQRe94m6eQ"
	ok = validator2.Validate(jwt2)
	if ok {
		t.Error(nil)
	}

	validator3 := jwt.JWTValidator{
		JWK:       "key1",
		Audiences: &audiences,
	}

	// has aud
	jwt3 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDc1NTc3OTMsIm5iZiI6MTcwNzU1Nzc5MywiYXVkIjpbImEiLCJjIiwiZCJdfQ.IIrinZQdXAQvRfv7MjYBqDfxNIqUZinoxsFrm9fziEw"
	ok = validator3.Validate(jwt3)
	if !ok {
		t.Error(nil)
	}

	// has no aud
	jwt3 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDc1NTc3OTMsIm5iZiI6MTcwNzU1Nzc5MywiYXVkIjpbImMiLCJkIl19.rilOQV1LJX0kUALR7S4p66t30UNMJLhe4m84824MaPc"
	ok = validator3.Validate(jwt3)
	if ok {
		t.Error(nil)
	}
}

func TestExtractJWTString(t *testing.T) {

	accessHeader := func(key string) (string, bool) {
		headers := map[string]string{"Auth": "XXX", "Authorization": "Bearer XXX"}
		value, ok := headers[key]
		return value, ok
	}

	accessQueryParam := func(key string) (string, bool) {
		params := map[string]string{"abc": "XXX"}
		value, ok := params[key]
		return value, ok
	}

	provider1 := jwt.JWTProvider{
		Name: "123",
		FromHeader: &jwt.JWTHeader{
			Name: "Auth",
		},
	}
	jwtString, err := provider1.ExtractJWTString(accessHeader, accessQueryParam)
	if err != nil || jwtString != "XXX" {
		t.Error(err)
	}

	fromParam := "abc"
	provider2 := jwt.JWTProvider{
		Name:      "123",
		FromParam: &fromParam,
	}
	jwtString, err = provider2.ExtractJWTString(accessHeader, accessQueryParam)
	if err != nil || jwtString != "XXX" {
		t.Error(err)
	}

	provider3 := jwt.JWTProvider{
		Name: "123",
	}
	jwtString, err = provider3.ExtractJWTString(accessHeader, accessQueryParam)
	if err != nil || jwtString != "XXX" {
		t.Error(err)
	}
}
