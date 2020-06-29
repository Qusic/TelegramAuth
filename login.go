package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func handleLogin(w http.ResponseWriter, r *http.Request, ctx *context) {
	valid, user := useToken(ctx.cookie, time.Now())
	if !valid {
		callback := fmt.Sprintf("%v/callback?%v=%v&%v=%v",
			config.pathPrefix,
			config.queryRole, url.QueryEscape(ctx.role),
			config.queryRedirect, url.QueryEscape(ctx.redirect))
		fmt.Fprint(w, loginPageStart)
		fmt.Fprintf(w, unauthenticatedBody, config.botName, callback)
		fmt.Fprint(w, loginPageEnd)
		return
	}
	access, ok := config.roleBindings[ctx.role][user]
	if !ok || !access {
		fmt.Fprint(w, loginPageStart)
		fmt.Fprintf(w, unauthorizedBody, user)
		fmt.Fprint(w, loginPageEnd)
		return
	}
	http.Redirect(w, r, ctx.redirect, http.StatusFound)
}

const loginPageStart = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
`

const loginPageEnd = `
</body>
</html>
`

const unauthenticatedBody = `
<script async src="https://telegram.org/js/telegram-widget.js"
	data-telegram-login="%v"
	data-size="large"
	data-auth-url="%v">
</script>
`

const unauthorizedBody = `
You have logged in as @%v, but are not authorized to access this resource.
`
