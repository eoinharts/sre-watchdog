package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	endpointUp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "watchdog_up",
			Help: "Whether the target endpoint is up: 1 for up, 0 for down",
		},
	)

	checksTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "watchdog_checks_total",
			Help: "Total number of endpoint checks performed",
		},
	)

	failuresTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "watchdog_failures_total",
			Help: "Total number of failed endpoint checks",
		},
	)

	requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "watchdog_request_duration_seconds",
			Help:    "Duration of endpoint checks in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(
		endpointUp,
		checksTotal,
		failuresTotal,
		requestDuration,
	)
}

func main() {
	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		targetURL = "https://example.com"
	}

	checkInterval := getCheckInterval()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	startMetricsServer()

	log.Printf(
		"watchdog started target=%s check_interval=%s metrics_port=8080",
		targetURL,
		checkInterval,
	)

	for {
		checkEndpoint(client, targetURL)
		time.Sleep(checkInterval)
	}
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Println("Prometheus metrics available at http://localhost:8080/metrics")

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Printf("metrics server failed: %v", err)
		}
	}()
}

func checkEndpoint(client *http.Client, targetURL string) bool {
	start := time.Now()
	checksTotal.Inc()

	resp, err := client.Get(targetURL)

	duration := time.Since(start)
	requestDuration.Observe(duration.Seconds())

	if err != nil {
		failuresTotal.Inc()
		endpointUp.Set(0)

		log.Printf(
			"target=%s up=false latency_ms=%d error=%q",
			targetURL,
			duration.Milliseconds(),
			err,
		)

		return false
	}

	defer resp.Body.Close()

	isUp := resp.StatusCode >= 200 && resp.StatusCode < 400

	if isUp {
		endpointUp.Set(1)
	} else {
		endpointUp.Set(0)
		failuresTotal.Inc()
	}

	fmt.Printf(
		"target=%s status=%d latency_ms=%d up=%t\n",
		targetURL,
		resp.StatusCode,
		duration.Milliseconds(),
		isUp,
	)

	return isUp
}

func getCheckInterval() time.Duration {
	interval := os.Getenv("CHECK_INTERVAL")

	if interval == "" {
		return 30 * time.Second
	}

	seconds, err := strconv.Atoi(interval)
	if err != nil || seconds <= 0 {
		log.Printf(
			"invalid CHECK_INTERVAL=%q, using default 30 seconds",
			interval,
		)
		return 30 * time.Second
	}

	return time.Duration(seconds) * time.Second
}
