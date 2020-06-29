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
		invalidTokenNoSignature, invalidTokenBadSignature := setupTestToken()
	for _, data := range []struct {
		binding string
		token   string
		valid   bool
	}{
		{binding: validUser, token: validToken1, valid: true},
		{binding: validUser, token: validToken2, valid: true},
		{binding: validUser, token: invalidTokenNoSignature, valid: false},
		{binding: validUser, token: invalidTokenBadSignature, valid: false},
		{binding: invalidUser, token: validToken1, valid: false},
		{binding: invalidUser, token: validToken2, valid: false},
		{binding: invalidUser, token: invalidTokenNoSignature, valid: false},
		{binding: invalidUser, token: invalidTokenBadSignature, valid: false},
	} {
		config.roleBindings = map[string]map[string]bool{role: {data.binding: true}}
		state.authCache = map[string]time.Time{}
		ctx := context{role: role, cookie: data.token}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
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
