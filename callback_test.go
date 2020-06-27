package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleCallback(t *testing.T) {
	config.cookieName = "cookie"
	config.cookiePath = "/"
	config.cookieDomain = "example.com"
	app := "test_app"
	_,
		_, _,
		validToken1, validToken2,
		invalidTokenNoSignature, invalidTokenBadSignature := setupTestToken()
	for _, data := range []struct {
		url   string
		token string
		valid bool
	}{
		{url: "/test_url_1", token: validToken1, valid: true},
		{url: "/test_url_2", token: validToken2, valid: true},
		{url: "/test_url_3", token: invalidTokenNoSignature, valid: false},
		{url: "/test_url_4", token: invalidTokenBadSignature, valid: false},
	} {
		config.appURL = map[string]string{app: data.url}
		state.authCache = map[string]time.Time{}
		ctx := context{app: app}
		rr := httptest.NewRecorder()
		url := url.URL{Path: "/", RawQuery: data.token}
		r := httptest.NewRequest(http.MethodGet, url.String(), nil)
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
			assert.Equal(t, data.url, location.String())
		} else {
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	}
}
