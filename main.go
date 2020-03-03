package main

import (
	"fmt"
	"time"
	"log"
	"os"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/kayamuskas/docker-registry-exporter/collector"
)

const defaultAddr = ":8080"

func main() {

    //Create a new instance of the collector and 
    //register it with the prometheus client.
    dirSizeCol := collector.NewDirSizeCol("/Users/kayama/Downloads", 5*time.Second)
    defer dirSizeCol.Stop()

    prometheus.MustRegister(dirSizeCol)

    addr := defaultAddr
    // $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("server starting to listen on %s", addr)
	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request: %s %s", r.Method, r.URL.Path)

		fmt.Fprintf(w, "Hello, world!")
	}))
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server listen error: %+v", err)
	}
}
