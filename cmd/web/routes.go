package main

import (
	"net/http"
	"time"
	h "github.com/heptiolabs/healthcheck"
		)

func (app *application) routes() *http.ServeMux {
  health := h.NewHandler()

	// Add a readiness check against the health of an upstream HTTP dependency
	health.AddReadinessCheck(
	    "upstream-dep-http",
	    h.HTTPGetCheck("0.0.0.0:4000", 500*time.Millisecond))

	health.AddLivenessCheck(
			"custom-check-timeout", h.GoroutineCountCheck(100))

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/healthz", health.ReadyEndpoint)

	// Sleep for just a moment to make sure our Async handler had a chance to run
	time.Sleep(500 * time.Millisecond)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}

