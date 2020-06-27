package main

import (
	"net/http"
	"time"
)

func handleCallback(w http.ResponseWriter, r *http.Request, ctx *context) {
	token := tokenFromQuery(r)
	valid, _ := useToken(token, time.Now())
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie := tokenToCookie(token)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, config.appURL[ctx.app], http.StatusFound)
}
