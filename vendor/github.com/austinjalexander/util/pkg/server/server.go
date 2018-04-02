package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const (
	timeout = time.Second * 20
)

// Config represents modifiable server configurations.
type Config struct {
	OnlyJSONresponses bool
}

// Handler represents a HTTP handler for the server router.
type Handler struct {
	Func                          func(http.ResponseWriter, *http.Request)
	Path                          string
	Headers, Methods, QueryParams []string
}

// JSONresponse reprents a JSON response body returned by the server.
type JSONresponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

// Run creates a new routed server and runs it.
func Run(cfg Config, handlers []Handler, port uint16) {
	r := configureMiddleware(cfg)

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Func).
			Headers(h.Headers...).
			Methods(h.Methods...).
			Queries(h.QueryParams...)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		Handler:      r,
		IdleTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("info: gracefully shutting down server...")
	os.Exit(0)
}

func configureMiddleware(cfg Config) *mux.Router {
	r := mux.NewRouter()
	if cfg.OnlyJSONresponses {
		r.Use(onlyJSONresponses)
	}
	return r
}

func onlyJSONresponses(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
