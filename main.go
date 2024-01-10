package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			temperature.Set(25.0)
			time.Sleep(recordInterval)
		}
	}()
}

const (
	// TODO: make this configurable
	recordInterval = 1 * time.Second
)

var (
	namespace   = "flat35hd99_private"
	temperature = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "temperature_in_celsius",
		Help:      "Current temperature in celsius",
	})
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":17818", nil)
}
