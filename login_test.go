package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleLogin(t *testing.T) {
	config.botName = "test_bot"
	config.pathPrefix = "test_prefix"
	config.cookieName = "cookie"
	app := "test_app"
	url := "/test_app_url"
	_,
		validUser, invalidUser,
		validToken1, validToken2,
		invalidTokenNoSignature, invalidTokenBadSignature := setupTestToken()
	config.appURL = map[string]string{app: url}
	for _, data := range []struct {
		access string
		token  string
		valid  int
	}{
		{access: validUser, token: validToken1, valid: 1},
		{access: validUser, token: validToken2, valid: 1},
		{access: validUser, token: invalidTokenNoSignature, valid: -1},
		{access: validUser, token: invalidTokenBadSignature, valid: -1},
		{access: invalidUser, token: validToken1, valid: 0},
		{access: invalidUser, token: validToken2, valid: 0},
		{access: invalidUser, token: invalidTokenNoSignature, valid: -1},
		{access: invalidUser, token: invalidTokenBadSignature, valid: -1},
	} {
		config.appAccess = map[string]map[string]bool{app: {data.access: true}}
		state.authCache = map[string]time.Time{}
		ctx := context{app: app}
		rr := httptest.NewRecorder()
		cookie := http.Cookie{Name: config.cookieName, Value: data.token}
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&cookie)
		handleLogin(rr, r, &ctx)
		switch data.valid {
		case 1:
			assert.Equal(t, http.StatusFound, rr.Code)
			location, _ := rr.Result().Location()
			assert.Equal(t, url, location.String())
		case 0:
			assert.Equal(t, http.StatusOK, rr.Code)
			body := rr.Body.String()
			assert.Contains(t, body, fmt.Sprintf("@%v", validUser))
		case -1:
			assert.Equal(t, http.StatusOK, rr.Code)
			body := rr.Body.String()
			assert.Contains(t, body, fmt.Sprintf("data-telegram-login=\"%v\"", config.botName))
			assert.Contains(t, body, fmt.Sprintf("data-auth-url=\"%v/%v/callback\"", config.pathPrefix, app))
		}
	}
}
