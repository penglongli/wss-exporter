package scheduler

import (
	"time"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"crypto/tls"
	"net/url"
)

var (
	urlStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "url_status_exporter",
			Help: "url response status_code exporter",
		},
		[]string{"url"},
	)
	urlTimeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "url_time_exporter",
			Help: "url request time exporter",
		},
		[]string{"url"},
	)
	client *http.Client
)

func init() {
	// Register gauge
	prometheus.MustRegister(urlStatusGauge)
	prometheus.MustRegister(urlTimeGauge)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
}

func CheckUrlStatus(slice []string, interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Second)
	for range t.C {
		for _, urlStr := range slice {
			go func(urlStr string) {
				check(urlStr)
			}(urlStr)
		}
	}
}

func check(urlStr string) {
	URL, err := url.Parse(urlStr)
	if err != nil {
		log.Errorf("Error format url: %s\n, errMsg: %s", urlStr, err.Error())
		return
	}

	switch scheme := URL.Scheme; scheme {
	case "http":
	case "https":
		checkHttpUrl(urlStr)
	case "ws":
		checkWsUrl("http", URL)
	case "wss":
		checkWsUrl("https", URL)
	default:
		log.Errorf("Unknown url protocol: %s\n", urlStr)
		return
	}
}

func checkWsUrl(newScheme string, urlPtr *url.URL) {
	withoutScheme := generateUrlWithoutScheme(urlPtr)
	urlStr := urlPtr.String()
	newUrlStr := newScheme + "://" + withoutScheme

	req, err := http.NewRequest("GET", newUrlStr, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Extensions", "permessage-deflate; client_max_window_bits")
	req.Header.Set("Sec-WebSocket-Protocol", "binary")
	req.Header.Set("Sec-WebSocket-Version", "13")
	// Just a random num, no actual meaning
	req.Header.Set("Sec-WebSocket-Key", "VJZpSmTxXtf5KMfh3MXe4w==")

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		urlStatusGauge.WithLabelValues(urlStr).Set(-1)
		return
	}
	urlStatusGauge.WithLabelValues(urlStr).Set(float64(resp.StatusCode))

	defer func() {
		resp.Body.Close()
		urlTimeGauge.WithLabelValues(urlStr).Set(time.Since(start).Seconds() * 1000)
	}()
}

func checkHttpUrl(urlStr string) {
	req, err := http.NewRequest("GET", urlStr, nil)

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		urlStatusGauge.WithLabelValues(urlStr).Set(-1)
		return
	}

	defer func() {
		resp.Body.Close()
		urlTimeGauge.WithLabelValues(urlStr).Set(time.Since(start).Seconds() * 1000)
	}()
	urlStatusGauge.WithLabelValues(urlStr).Set(float64(resp.StatusCode))
}

func generateUrlWithoutScheme(urlPtr *url.URL) string {
	withoutScheme := urlPtr.Host + urlPtr.Path

	rawQuery := urlPtr.RawQuery
	if rawQuery == "" {
		return withoutScheme
	}
	withoutScheme = withoutScheme + "?" + urlPtr.RawQuery
	return withoutScheme
}
