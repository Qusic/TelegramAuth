package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleLogin(t *testing.T) {
	config.botName = "test_bot"
	config.pathPrefix = "test_prefix"
	role := "test_role"
	redirect := "https://example.com/test_redirect"
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
		valid   int
	}{
		{binding: validUser, token: validToken1, valid: 1},
		{binding: validUser, token: validToken2, valid: 1},
		{binding: validUser, token: invalidTokenNoSignature, valid: -1},
		{binding: validUser, token: invalidTokenBadSignature, valid: -1},
		{binding: validUser, token: invalidTokenBadString, valid: -1},
		{binding: validUser, token: invalidTokenBadTimestamp, valid: -1},
		{binding: invalidUser, token: validToken1, valid: 0},
		{binding: invalidUser, token: validToken2, valid: 0},
		{binding: invalidUser, token: invalidTokenNoSignature, valid: -1},
		{binding: invalidUser, token: invalidTokenBadSignature, valid: -1},
		{binding: invalidUser, token: invalidTokenBadString, valid: -1},
		{binding: invalidUser, token: invalidTokenBadTimestamp, valid: -1},
	} {
		config.roleBindings = map[string]map[string]bool{role: {data.binding: true}}
		state.authCache = map[string]time.Time{}
		ctx := context{role: role, redirect: redirect, cookie: data.token}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		handleLogin(rr, r, &ctx)
		switch data.valid {
		case 1:
			assert.Equal(t, http.StatusFound, rr.Code)
			location, _ := rr.Result().Location()
			assert.Equal(t, redirect, location.String())
		case 0:
			assert.Equal(t, http.StatusOK, rr.Code)
			body := rr.Body.String()
			assert.Contains(t, body, "@"+validUser)
		case -1:
			assert.Equal(t, http.StatusOK, rr.Code)
			body := rr.Body.String()
			assert.Contains(t, body, config.botName)
			assert.Contains(t, body, config.pathPrefix+"/callback")
			assert.Contains(t, body, config.queryRole+"="+url.QueryEscape(role))
			assert.Contains(t, body, config.queryRedirect+"="+url.QueryEscape(redirect))
		}
	}
}
