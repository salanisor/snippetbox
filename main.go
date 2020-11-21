/*
Copyright 2020 Red Hat

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/klauspost/cpuid"
        h "github.com/heptiolabs/healthcheck"
)

var GitCommit string
var Arch string
var Built string
var GoVersion string

// health return
func healthy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API accepting connections"))
}

// version returns the api version
func version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("{Git Commit:\"%s\",CPU_arch:%s,Built:\"%s\",Go_version:%s}\n", GitCommit, Arch, Built, GoVersion)))
}

// hostCPU returns the cpu info
func hostCPU(w http.ResponseWriter, r *http.Request) {
	w.Write(
		[]byte(fmt.Sprintf("{name:\"%s\",model:\"%d\",family:\"%d\"}\n", cpuid.CPU.BrandName, cpuid.CPU.Model, cpuid.CPU.Family)),
	)
}

// health of the app
func health(w http.ResponseWriter, r *http.Request) {
	// Create a Handler that we can use to register liveness and readiness checks.
	health := h.NewHandler()

	// Add a readiness check against the health of an upstream HTTP dependency
	upstreamURL := "http://0.0.0.0:8080/healthy"
	health.AddReadinessCheck(
		"upstream-dep-http",
		h.HTTPGetCheck(upstreamURL, 500*time.Millisecond))

	// Implement a custom check with a 50 millisecond timeout.
	health.AddLivenessCheck("custom-check-timeout", h.GoroutineCountCheck(100))

	// Expose the readiness endpoints on a custom path /healthz mixed into
	// our main application mux.
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.ReadyEndpoint)

	// Sleep for just a moment to make sure our Async handler had a chance to run
	time.Sleep(500 * time.Millisecond)
}
func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/healthy", healthy)
	r.HandleFunc("/version", version)
	r.HandleFunc("/cpu", hostCPU)
	r.HandleFunc("/healthz", health)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
