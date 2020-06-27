package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	initialize("config.yaml", log.Fatalln)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	r.Route(config.pathPrefix+"/{app}", func(r chi.Router) {
		r.Get("/", useContext(handleAuth))
		r.Get("/login", useContext(handleLogin))
		r.Get("/callback", useContext(handleCallback))
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	log.Printf(startMessage,
		config.botName, config.listenAddress, config.pathPrefix,
		config.cookieName, config.authHeader,
		config.authDuration.String(), config.authTimeout.String())
	err := http.ListenAndServe(config.listenAddress, r)
	if err != nil {
		log.Fatalln("fail to start server:", err)
	}
}

const startMessage = `
TelegramAuth is starting for @%v on %v%v
Cookie name for clients is "%v". Header name for servers is "%v".
New tokens must be created in the past [%v] to be accepted.
Accepted tokens can outlive the lifespan unless being unused in the past [%v].
`
