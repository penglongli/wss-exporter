package scheduler

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func CheckWssStatus(url string, interval int, gauge *prometheus.GaugeVec) {
	t := time.NewTimer(time.Duration(interval) * time.Second)
	for range t.C {
		check(url, gauge)
	}
}

func check(url string, gauge *prometheus.GaugeVec) {
	gauge.WithLabelValues("lipenglong").Set(1)
}