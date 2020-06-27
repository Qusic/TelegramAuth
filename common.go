package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

type context struct {
	app string
}

var (
	config struct {
		botName       string
		botToken      string
		listenAddress string
		pathPrefix    string
		cookieName    string
		cookiePath    string
		cookieDomain  string
		authHeader    string
		authDuration  time.Duration
		authTimeout   time.Duration
		appURL        map[string]string
		appAccess     map[string]map[string]bool
	}
	state struct {
		authCache map[string]time.Time
		authMutex sync.Mutex
	}
)

func initialize(file string, fatal func(v ...interface{})) {
	type configApp struct {
		Name  string
		URL   string
		Users []string
	}
	type configRoot struct {
		Global struct {
			Bot     string
			Token   string
			Address string
			Prefix  string
		}
		Cookie struct {
			Name   string
			Path   string
			Domain string
		}
		Auth struct {
			Header   string
			Duration string
			Timeout  string
		}
		Apps []configApp
	}
	readString := func(str string, field string, required bool) string {
		if required && str == "" {
			fatal("bad config field "+field+":", "required but missing")
		}
		return str
	}
	readStrings := func(strs []string, field string, required bool) []string {
		if required && len(strs) == 0 {
			fatal("bad config field "+field+":", "required but missing")
		}
		return strs
	}
	readDuration := func(str, field string, required bool) time.Duration {
		str = readString(str, field, required)
		duration, err := time.ParseDuration(str)
		if required && err != nil {
			fatal("bad config field "+field+":", err)
		}
		return duration
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fatal("fail to read config file:", err)
	}
	root := configRoot{}
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		fatal("fail to unmarshal config file:", err)
	}
	config.botName = readString(root.Global.Bot, "global.bot", true)
	config.botToken = readString(root.Global.Token, "global.token", true)
	config.listenAddress = readString(root.Global.Address, "global.address", true)
	config.pathPrefix = readString(root.Global.Prefix, "global.prefix", false)
	config.cookieName = readString(root.Cookie.Name, "cookie.name", true)
	config.cookiePath = readString(root.Cookie.Path, "cookie.path", false)
	config.cookieDomain = readString(root.Cookie.Domain, "cookie.domain", false)
	config.authHeader = readString(root.Auth.Header, "auth.header", false)
	config.authDuration = readDuration(root.Auth.Duration, "auth.duration", true)
	config.authTimeout = readDuration(root.Auth.Timeout, "auth.timeout", true)
	config.appURL = map[string]string{}
	config.appAccess = map[string]map[string]bool{}
	for index, app := range root.Apps {
		field := fmt.Sprintf("apps[%v].", index)
		name := readString(app.Name, field+"name", true)
		url := readString(app.URL, field+"url", true)
		users := readStrings(app.Users, field+"users", true)
		config.appURL[name] = url
		config.appAccess[name] = map[string]bool{}
		for _, user := range users {
			config.appAccess[name][user] = true
		}
	}
	state.authCache = map[string]time.Time{}
}

func useContext(next func(http.ResponseWriter, *http.Request, *context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context{}
		ctx.app = chi.URLParam(r, "app")
		if _, ok := config.appURL[ctx.app]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next(w, r, &ctx)
	}
}

func useToken(token string, now time.Time) (valid bool, user string) {
	data, err := url.ParseQuery(token)
	if err != nil {
		return
	}
	signature := data.Get("hash")
	data.Del("hash")
	message, _ := url.QueryUnescape(strings.ReplaceAll(data.Encode(), "&", "\n"))
	key := sha256.Sum256([]byte(config.botToken))
	mac := hmac.New(sha256.New, key[:])
	io.WriteString(mac, message)
	if hex.EncodeToString(mac.Sum(nil)) != signature {
		return
	}
	timestamp, err := strconv.ParseInt(data.Get("auth_date"), 10, 64)
	if err != nil {
		return
	}
	createdTime := time.Unix(timestamp, 0)
	recentlyCreated := now.Sub(createdTime) < config.authDuration
	state.authMutex.Lock()
	defer state.authMutex.Unlock()
	usedTime, ok := state.authCache[token]
	recentlyUsed := ok && now.Sub(usedTime) < config.authTimeout
	if recentlyCreated || recentlyUsed {
		state.authCache[token] = now
		valid = true
		user = data.Get("username")
	} else {
		delete(state.authCache, token)
	}
	return
}

func tokenFromQuery(r *http.Request) (token string) {
	token = r.URL.RawQuery
	return token
}

func tokenFromCookie(r *http.Request) (token string) {
	cookie, err := r.Cookie(config.cookieName)
	if err != nil {
		return
	}
	token = cookie.Value
	return
}

func tokenToCookie(token string) (cookie http.Cookie) {
	cookie.Name = config.cookieName
	cookie.Value = token
	cookie.Path = config.cookiePath
	cookie.Domain = config.cookieDomain
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	return
}
