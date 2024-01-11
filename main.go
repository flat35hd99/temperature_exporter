package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
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
				"id": "sample1",
			}).Set(25.0 + rand.Float64()*5.0 - 2.5)
			temperature.With(prometheus.Labels{
				"id": "sample2",
			}).Set(25.0 + rand.Float64()*5.0 - 2.5)

			// List directories of /sys/bus/w1/devices
			dirEntries, err := os.ReadDir("/sys/bus/w1/devices")
			if err != nil {
				time.Sleep(recordInterval)
				continue
			}
			for _, dirEntry := range dirEntries {
				if !dirEntry.IsDir() {
					continue
				}
				deviceId := dirEntry.Name()

				buf, err := os.ReadFile("/sys/bus/w1/devices/" + deviceId + "/w1_slave")
				if err != nil {
					continue
				}
				regexp := regexp.MustCompile(`t=(-?\d+)`)
				matches := regexp.FindAllSubmatch(buf, -1)
				for _, match := range matches {
					// Convert types: byte -> string -> int -> float64
					intTemperatureInMilliCelsius, err := strconv.Atoi(string(match[1]))
					if err != nil {
						continue
					}
					temperatureInCelsius := float64(intTemperatureInMilliCelsius) / 1000.0

					temperature.With(prometheus.Labels{
						"id": deviceId,
					}).Set(temperatureInCelsius)
				}
			}
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
