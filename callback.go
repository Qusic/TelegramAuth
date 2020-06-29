package main

import (
	"net/http"
	"time"
)

func handleCallback(w http.ResponseWriter, r *http.Request, ctx *context) {
	valid, _ := useToken(ctx.query, time.Now())
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie := buildCookie(ctx.query)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, ctx.redirect, http.StatusFound)
}
