package main

import (
    "fmt"
    "time"
    "log"
    "os"
    "net/http"
    "path/filepath"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

const defaultAddr = ":8080"
var metricValue float64

// home logs the received request and returns a simple response.
func home(w http.ResponseWriter, r *http.Request) {
    log.Printf("received request: %s %s", r.Method, r.URL.Path)
    fmt.Fprintf(w, "Hello, world!")
}

func DirSize(path string) (int64, error){
    var size int64
    adjSize := func(_ string, info os.FileInfo, err error) error {
	if err != nil {
	    return err
	}
	if !info.IsDir() {
	    size += info.Size()
	}
	return err
    }

    err := filepath.Walk(path, adjSize)    
    return size, err

}

//First run to get initial values
func init() {
    fmt.Println(DirSize("/Users/kayama/Downloads"))
}

func main() {

    //Create a new instance of the collector and 
    //register it with the prometheus client.
    run := NewFilesystemCollector()
    prometheus.MustRegister(run)

    ticker := time.NewTicker(5 * time.Second)

    go func() {
        for {
            select {
        	case t := <-ticker.C:
                fmt.Println("Tick at", t)
                //metricValue := fmt.Println(DirSize("/Users/kayama/Downloads"))
                fmt.Println(DirSize("/Users/kayama/Downloads"))
                var mmm int64 = 1001
                var zzz float64 = float64(mmm) + 0.001
                fmt.Println(zzz)
                
                metricValue = 100.1
                fmt.Println(metricValue)
            }
        }
    }()

    addr := defaultAddr
    // $PORT environment variable is provided in the Kubernetes deployment.
    if p := os.Getenv("PORT"); p != "" {
	addr = ":" + p
    }

    log.Printf("server starting to listen on %s", addr)
    http.HandleFunc("/", home)
    http.Handle("/metrics", promhttp.Handler())

    if err := http.ListenAndServe(addr, nil); err != nil {
	log.Fatalf("server listen error: %+v", err)
    }

}
