package main

import (
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"wss-exporter/scheduler"
)

var (
	wssStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wss_status",
			Help: "wss status exporter",
		},
		[]string{"wss_url"},
	)
	cfgFile      string
	ListenPort   string
	WssUrl       string
	TimeInterval int
)

func init() {
	prometheus.MustRegister(wssStatus)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		log.Info("Using config file: " + viper.ConfigFileUsed())
	}

	// Init some variables
	initVars()
}

func initVars() {
	ListenPort = viper.GetString("port")
    if ListenPort == "" {
        ListenPort = ":8080"
    }

	WssUrl = viper.GetString("wss_url")
	TimeInterval = viper.GetInt("time_interval")
}

func main() {
	go scheduler.CheckWssStatus(WssUrl, TimeInterval, wssStatus)
	http.Handle("/metrics", promhttp.Handler())

	// Start http server
	err := http.ListenAndServe(ListenPort, nil)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// Block util a signal is received
	log.Println(<-ch)
}
