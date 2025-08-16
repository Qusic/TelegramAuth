package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestToken() (
	now time.Time,
	validUser, invalidUser string,
	validToken1, validToken2 string,
	invalidTokenNoSignature string,
	invalidTokenBadSignature string,
	invalidTokenBadString string,
	invalidTokenBadTimestamp string,
) {
	config.botName = "test_bot"
	config.botToken = "test_bot_fake_token"
	config.authDuration = time.Second
	config.authTimeout = time.Second
	stringify := func(data map[string]string, escaped, sorted bool, separator string) string {
		fields := make([]string, 0, len(data))
		for key, value := range data {
			field := key + "="
			if escaped {
				field += url.QueryEscape(value)
			} else {
				field += value
			}
			fields = append(fields, field)
		}
		if sorted {
			sort.Strings(fields)
		}
		return strings.Join(fields, separator)
	}
	sign := func(s string) string {
		key := sha256.Sum256([]byte(config.botToken))
		mac := hmac.New(sha256.New, key[:])
		io.WriteString(mac, s)
		return hex.EncodeToString(mac.Sum(nil))
	}
	now = time.Now()
	validUser = "test_user"
	invalidUser = ""
	message := map[string]string{
		"id":         "1234",
		"first_name": "名",
		"last_name":  "姓",
		"username":   validUser,
		"auth_date":  strconv.FormatInt(now.Unix(), 10),
	}
	signature := sign(stringify(message, false, true, "\n"))
	sortedMessage := stringify(message, true, true, "&")
	unsortedMessage := stringify(message, true, false, "&")
	validToken1 = sortedMessage + "&hash=" + signature
	validToken2 = unsortedMessage + "&hash=" + signature
	invalidTokenNoSignature = sortedMessage
	invalidTokenBadSignature = validToken1 + "0"
	invalidTokenBadString = "%xx"
	message["auth_date"] = "xxxx"
	invalidTokenBadTimestamp = fmt.Sprintf("%v&hash=%v",
		stringify(message, true, true, "&"),
		sign(stringify(message, false, true, "\n")))
	return
}

func TestInitialize(t *testing.T) {
	initialize("config.example.yml", t.Fatal)
	assert.Equal(t, "MyTelegramBot", config.botName)
	assert.Equal(t, "1234:ABCDEFG", config.botToken)
	assert.Equal(t, "tcp", config.listenNetwork)
	assert.Equal(t, ":80", config.listenAddress)
	assert.Equal(t, "/auth", config.pathPrefix)
	assert.Equal(t, "role", config.queryRole)
	assert.Equal(t, "redirect_uri", config.queryRedirect)
	assert.Equal(t, "token", config.cookieName)
	assert.Equal(t, "/", config.cookiePath)
	assert.Equal(t, "example.com", config.cookieDomain)
	assert.Equal(t, "X-Telegram-Auth", config.authHeader)
	assert.Equal(t, 12*time.Hour, config.authDuration)
	assert.Equal(t, 10*time.Minute, config.authTimeout)
	assert.Equal(t, map[string]map[string]bool{
		"owner": {
			"UserMe": true,
		},
		"contributor": {
			"UserMe": true,
			"UserA":  true,
		},
		"123": {
			"UserA": true,
			"UserB": true,
			"UserC": true,
		},
		"0": {
			"UserD": true,
			"UserE": true,
			"UserF": true,
		},
	}, config.roleBindings)
	assert.Equal(t, map[string]time.Time{}, state.authCache)
}

func TestUseContext(t *testing.T) {
	config.queryRole = "roletest"
	config.queryRedirect = "redirecttest"
	config.cookieName = "cookietest"
	config.roleBindings = map[string]map[string]bool{
		"rb": {},
	}
	for _, data := range []struct {
		valid   bool
		query   string
		cookie  string
		context context
	}{
		{valid: true,
			query:  "roletest=rb&redirecttest=https%3A%2F%2Fexample.com%2Ftest&test=111",
			cookie: "cookietest=123",
			context: context{
				role: "rb", redirect: "https://example.com/test",
				query: "test=111", cookie: "123",
			},
		},
		{valid: true,
			query:  "roletest=rb&redirecttest=https%3A%2F%2Fexample.com%2Ftest&test=111",
			cookie: "cookie=123",
			context: context{
				role: "rb", redirect: "https://example.com/test",
				query: "test=111", cookie: "",
			},
		},
		{valid: true,
			query:  "roletest=rb&test=111",
			cookie: "",
			context: context{
				role: "rb", redirect: "",
				query: "test=111", cookie: "",
			},
		},
		{valid: true,
			query:  "roletest=rb",
			cookie: "cookietest=123",
			context: context{
				role: "rb", redirect: "",
				query: "", cookie: "123",
			},
		},
		{valid: true,
			query:  "roletest=rb",
			cookie: "",
			context: context{
				role: "rb", redirect: "",
				query: "", cookie: "",
			},
		},
		{valid: false,
			query:  "roletest=rb123",
			cookie: "cookietest=123",
			context: context{
				role: "", redirect: "",
				query: "", cookie: "",
			},
		},
		{valid: false,
			query:  "redirecttest=https%3A%2F%2Fexample.com%2Ftest",
			cookie: "cookietest=123",
			context: context{
				role: "", redirect: "",
				query: "", cookie: "",
			},
		},
		{valid: false,
			query:  "",
			cookie: "",
			context: context{
				role: "", redirect: "",
				query: "", cookie: "",
			},
		},
	} {
		called := false
		ctx := context{}
		handler := func(w http.ResponseWriter, r *http.Request, c *context) {
			called = true
			ctx = *c
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?"+data.query, nil)
		r.Header.Set("Cookie", data.cookie)
		useContext(handler).ServeHTTP(rr, r)
		if data.valid {
			assert.True(t, called)
			assert.Equal(t, data.context, ctx)
			assert.Equal(t, http.StatusOK, rr.Code)
		} else {
			assert.False(t, called)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	}
}

func TestUseToken(t *testing.T) {
	now,
		validUser, invalidUser,
		validToken1, validToken2,
		invalidTokenNoSignature,
		invalidTokenBadSignature,
		invalidTokenBadString,
		invalidTokenBadTimestamp := setupTestToken()
	for _, data := range []struct {
		token string
		valid bool
	}{
		{token: validToken1, valid: true},
		{token: validToken2, valid: true},
		{token: invalidTokenNoSignature, valid: false},
		{token: invalidTokenBadSignature, valid: false},
		{token: invalidTokenBadString, valid: false},
		{token: invalidTokenBadTimestamp, valid: false},
	} {
		state.authCache = map[string]time.Time{}
		valid, user := useToken(data.token, now)
		if data.valid {
			assert.True(t, valid)
			assert.Equal(t, validUser, user)
		} else {
			assert.False(t, valid)
			assert.Equal(t, invalidUser, user)
		}
	}
	for _, data := range []struct {
		duration int
		timeout  int
		offsets  []int
		valids   []bool
	}{
		{duration: 0, timeout: 0,
			offsets: []int{0, 2, 4, 8},
			valids:  []bool{false, false, false, false},
		},
		{duration: 10, timeout: 0,
			offsets: []int{8, 9, 10, 11},
			valids:  []bool{true, true, false, false},
		},
		{duration: 100, timeout: 5,
			offsets: []int{10, 90, 100, 110},
			valids:  []bool{true, true, false, false},
		},
		{duration: 100, timeout: 5,
			offsets: []int{99, 102, 105, 110},
			valids:  []bool{true, true, true, false},
		},
		{duration: 100, timeout: 5,
			offsets: []int{100, 101, 102, 103},
			valids:  []bool{false, false, false, false},
		},
	} {
		config.authDuration = time.Duration(data.duration) * time.Second
		config.authTimeout = time.Duration(data.timeout) * time.Second
		state.authCache = map[string]time.Time{}
		for index, offset := range data.offsets {
			now := now.Add(time.Duration(offset) * time.Second)
			valid, user := useToken(validToken1, now)
			if data.valids[index] {
				assert.True(t, valid)
				assert.Equal(t, validUser, user)
				assert.Equal(t, map[string]time.Time{validToken1: now}, state.authCache)
			} else {
				assert.False(t, valid)
				assert.Equal(t, invalidUser, user)
				assert.Equal(t, map[string]time.Time{}, state.authCache)
			}
		}
	}
}
