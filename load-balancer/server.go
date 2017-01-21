package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "requests",
	Help: "Number of requests",
})

type handler struct {
	latency time.Duration
}

func (h *handler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	requests.Inc()
	time.Sleep(h.latency)
	w.Write([]byte("pong"))
}

func init() {
	prometheus.MustRegister(requests)
}

func main() {
	addr := flag.String("addr", ":8501", "service port to run on")
	latency := flag.Duration("latency", time.Duration(0), "latency to add to each request")
	flag.Parse()

	fmt.Printf("serving on %s with %v latency\n", *addr, *latency)

	httpHandler := handler{latency: *latency}
	http.HandleFunc("/", httpHandler.HandleRequest)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*addr, nil)
}
