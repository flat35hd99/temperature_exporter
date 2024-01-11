package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			temperature.With(prometheus.Labels{
				"id": "asdf",
			}).Set(25.0 + rand.Float64()*5.0 - 2.5)
			temperature.With(prometheus.Labels{
				"id": "qwer",
			}).Set(25.0 + rand.Float64()*5.0 - 2.5)
			time.Sleep(recordInterval)
		}
	}()
}

const (
	// TODO: make this configurable
	recordInterval = 1 * time.Second
	port           = 17818
)

var (
	namespace   = "flat35hd99_private"
	temperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "temperature_in_celsius",
		Help:      "Current temperature in celsius",
	}, []string{
		"id",
	})
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	log.Default().Printf("Starting server on port %d", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
