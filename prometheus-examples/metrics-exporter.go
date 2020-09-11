package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Way to check whether the belwo metrics 
// have been registered and being scraped

/*
	[root@k8s-node1 prometheus-2.18.1.linux-amd64]# curl http://k8s-master:4321/metrics | grep hd_errors
	# HELP hd_errors_total Number of hard-disk errors.
	# TYPE hd_errors_total counter
	hd_errors_total{device="/dev/sda"} 6

*/

var (
	// examples for counter and gauge metrics
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})

	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
)

// Metrics have to registered first,
// these can be either called in init() or called from main function
// before setting the values.
func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)

}

func main() {
	cpuTemp.Set(65.3)

	go func() {
		hdFailures.WithLabelValues("/dev/sda").Inc()
		time.Sleep(10*time.Second)
	} ()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":4321", nil))
}
