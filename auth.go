package main

import (
	"net/http"
	"time"
)

func handleAuth(w http.ResponseWriter, r *http.Request, ctx *context) {
	valid, user := useToken(ctx.cookie, time.Now())
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	access, ok := config.roleBindings[ctx.role][user]
	if !ok || !access {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if config.authHeader != "" {
		w.Header().Set(config.authHeader, user)
	}
	w.WriteHeader(http.StatusOK)
}
