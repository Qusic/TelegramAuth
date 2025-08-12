package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	initialize("config.yaml", log.Fatalln)
	go runTask()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	r.Route(config.pathPrefix, func(r chi.Router) {
		r.Handle("/", useContext(handleAuth))
		r.Handle("/login", useContext(handleLogin))
		r.Handle("/callback", useContext(handleCallback))
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	log.Printf("TelegramAuth is starting for @%v on %v%v",
		config.botName, config.listenAddress, config.pathPrefix)
	err := http.ListenAndServe(config.listenAddress, r)
	if err != nil {
		log.Fatalln("fail to start server:", err)
	}
}
