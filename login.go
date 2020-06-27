package main

import (
	"fmt"
	"net/http"
)

func handleLogin(w http.ResponseWriter, r *http.Request, ctx *context) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, loginPage, config.botName, config.pathPrefix, ctx.app)
}

const loginPage = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<style>
			body {
				width: 100vw;
				height: 100vh;
				margin: 0;
				display: flex;
				align-items: center;
				justify-content: center;
			}
		</style>
	</head>
	<body>
		<script async src="https://telegram.org/js/telegram-widget.js"
			data-telegram-login="%v"
			data-size="large"
			data-auth-url="%v/%v/callback">
		</script>
	</body>
</html>
`
