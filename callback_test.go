package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleCallback(t *testing.T) {
	config.cookieName = "cookie"
	config.cookiePath = "/"
	config.cookieDomain = "example.com"
	role := "test_role"
	_,
		_, _,
		validToken1, validToken2,
		invalidTokenNoSignature,
		invalidTokenBadSignature,
		invalidTokenBadString,
		invalidTokenBadTimestamp := setupTestToken()
	for _, data := range []struct {
		redirect string
		token    string
		valid    bool
	}{
		{redirect: "/test_redirect_1", token: validToken1, valid: true},
		{redirect: "/test_redirect_2", token: validToken2, valid: true},
		{redirect: "/test_redirect_3", token: invalidTokenNoSignature, valid: false},
		{redirect: "/test_redirect_4", token: invalidTokenBadSignature, valid: false},
		{redirect: "/test_redirect_5", token: invalidTokenBadString, valid: false},
		{redirect: "/test_redirect_6", token: invalidTokenBadTimestamp, valid: false},
	} {
		state.authCache = map[string]time.Time{}
		ctx := context{role: role, redirect: data.redirect, query: data.token}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		handleCallback(rr, r, &ctx)
		if data.valid {
			assert.Equal(t, http.StatusFound, rr.Code)
			cookies := rr.Result().Cookies()
			assert.Len(t, cookies, 1)
			assert.Equal(t, config.cookieName, cookies[0].Name)
			assert.Equal(t, data.token, cookies[0].Value)
			assert.Equal(t, config.cookiePath, cookies[0].Path)
			assert.Equal(t, config.cookieDomain, cookies[0].Domain)
			location, _ := rr.Result().Location()
			assert.Equal(t, data.redirect, location.String())
		} else {
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	}
}
