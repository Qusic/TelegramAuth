package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleLogin(t *testing.T) {
	for _, data := range []struct {
		botName, pathPrefix, app string
	}{
		{botName: "bot", pathPrefix: "/", app: "app"},
		{botName: "test_bot", pathPrefix: "/auth", app: "test_app"},
	} {
		config.botName = data.botName
		config.pathPrefix = data.pathPrefix
		ctx := context{app: data.app}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		handleLogin(rr, r, &ctx)
		assert.Equal(t, http.StatusOK, rr.Code)
		body := rr.Body.String()
		assert.Contains(t, body, fmt.Sprintf("data-telegram-login=\"%v\"", data.botName))
		assert.Contains(t, body, fmt.Sprintf("data-auth-url=\"%v/%v/callback\"", data.pathPrefix, data.app))
	}
}
