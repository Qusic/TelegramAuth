package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleAuth(t *testing.T) {
	config.cookieName = "cookie"
	config.authHeader = "X-Test-Header"
	app := "test_app"
	_,
		validUser, invalidUser,
		validToken1, validToken2,
		invalidTokenNoSignature, invalidTokenBadSignature := setupTestToken()
	for _, data := range []struct {
		access string
		token  string
		valid  bool
	}{
		{access: validUser, token: validToken1, valid: true},
		{access: validUser, token: validToken2, valid: true},
		{access: validUser, token: invalidTokenNoSignature, valid: false},
		{access: validUser, token: invalidTokenBadSignature, valid: false},
		{access: invalidUser, token: validToken1, valid: false},
		{access: invalidUser, token: validToken2, valid: false},
		{access: invalidUser, token: invalidTokenNoSignature, valid: false},
		{access: invalidUser, token: invalidTokenBadSignature, valid: false},
	} {
		config.appAccess = map[string]map[string]bool{app: {data.access: true}}
		state.authCache = map[string]time.Time{}
		ctx := context{app: app}
		rr := httptest.NewRecorder()
		cookie := http.Cookie{Name: config.cookieName, Value: data.token}
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&cookie)
		handleAuth(rr, r, &ctx)
		if data.valid {
			assert.Equal(t, http.StatusOK, rr.Code)
			header, ok := rr.HeaderMap[config.authHeader]
			assert.True(t, ok)
			assert.Equal(t, []string{validUser}, header)
		} else {
			assert.Equal(t, http.StatusUnauthorized, rr.Code)
			header, ok := rr.HeaderMap[config.authHeader]
			assert.False(t, ok)
			assert.Nil(t, header)
		}
	}
}
