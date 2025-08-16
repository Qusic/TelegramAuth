package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	arg := "config.yml"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	initialize(arg, log.Fatalln)
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

	log.Printf("TelegramAuth is starting for @%v on %v://%v%v",
		config.botName, config.listenNetwork, config.listenAddress, config.pathPrefix)
	listener, err := net.Listen(config.listenNetwork, config.listenAddress)
	if err != nil {
		log.Fatalln("fail to create listener:", err)
	}
	err = http.Serve(listener, r)
	if err != nil {
		log.Fatalln("fail to serve http:", err)
	}
}
