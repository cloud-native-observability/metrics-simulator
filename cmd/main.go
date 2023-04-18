package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cloud-native-observability/metrics-simulator/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <config file abs path>\n", os.Args[0])
		os.Exit(1)
	}
	configFilePath := os.Args[1]
	registry := prometheus.NewRegistry()
	generator, err := metrics.NewGenerator(configFilePath, registry)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	generator.Load()
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8888", nil)
}
