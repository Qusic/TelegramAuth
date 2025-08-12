package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleAuth(t *testing.T) {
	config.authHeader = "X-Test-Header"
	role := "test_role"
	_,
		validUser, invalidUser,
		validToken1, validToken2,
		invalidTokenNoSignature,
		invalidTokenBadSignature,
		invalidTokenBadString,
		invalidTokenBadTimestamp := setupTestToken()
	for _, data := range []struct {
		binding string
		token   string
		valid   bool
	}{
		{binding: validUser, token: validToken1, valid: true},
		{binding: validUser, token: validToken2, valid: true},
		{binding: validUser, token: invalidTokenNoSignature, valid: false},
		{binding: validUser, token: invalidTokenBadSignature, valid: false},
		{binding: validUser, token: invalidTokenBadString, valid: false},
		{binding: validUser, token: invalidTokenBadTimestamp, valid: false},
		{binding: invalidUser, token: validToken1, valid: false},
		{binding: invalidUser, token: validToken2, valid: false},
		{binding: invalidUser, token: invalidTokenNoSignature, valid: false},
		{binding: invalidUser, token: invalidTokenBadSignature, valid: false},
		{binding: invalidUser, token: invalidTokenBadString, valid: false},
		{binding: invalidUser, token: invalidTokenBadTimestamp, valid: false},
	} {
		config.roleBindings = map[string]map[string]bool{role: {data.binding: true}}
		state.authCache = map[string]time.Time{}
		ctx := context{role: role, cookie: data.token}
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handleAuth(rr, req, &ctx)
		res := rr.Result()
		if data.valid {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			header, ok := res.Header[config.authHeader]
			assert.True(t, ok)
			assert.Equal(t, []string{validUser}, header)
		} else {
			assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
			header, ok := res.Header[config.authHeader]
			assert.False(t, ok)
			assert.Nil(t, header)
		}
	}
}
