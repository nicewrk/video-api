package main

import (
	"github.com/austinjalexander/util/pkg/server"
	"github.com/nicewrk/video-api/api/healthcheck"
)

func main() {
	apiDoc := apidoc.New()

	cfg := server.Config{
		OnlyJSONresponses: true,
	}
	handlers := []server.Handler{
		apiDoc.Write(healthcheck.Configure()),
	}
	server.Run(cfg, handlers, 8080)
}
