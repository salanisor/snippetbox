package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/klauspost/cpuid"
)

// "something"
var GitCommit string

// "something"
var Arch string

// "something"
var Built string

// "something"
var GoVersion string

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
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

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/cpu", hostCPU)
	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":8080")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
