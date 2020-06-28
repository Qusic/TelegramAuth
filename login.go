package main

import (
	"fmt"
	"net/http"
	"time"
)

func handleLogin(w http.ResponseWriter, r *http.Request, ctx *context) {
	token := tokenFromCookie(r)
	valid, user := useToken(token, time.Now())
	if !valid {
		fmt.Fprint(w, loginPageStart)
		fmt.Fprintf(w, unauthenticatedBody, config.botName, config.pathPrefix, ctx.app)
		fmt.Fprint(w, loginPageEnd)
		return
	}
	access, ok := config.appAccess[ctx.app][user]
	if !ok || !access {
		fmt.Fprint(w, loginPageStart)
		fmt.Fprintf(w, unauthorizedBody, user)
		fmt.Fprint(w, loginPageEnd)
		return
	}
	http.Redirect(w, r, config.appURL[ctx.app], http.StatusFound)
}

const loginPageStart = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
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
	data-auth-url="%v/%v/callback">
</script>
`

const unauthorizedBody = `
You have logged in as @%v, but are not authorized to access this resource.
`
