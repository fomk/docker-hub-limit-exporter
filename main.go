package main

import (
	"flag"
	"github.com/fomk/docker-hub-limit-exporter/collector"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	c "github.com/fomk/docker-hub-limit-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	BuildVersion = "Development version"
	BuildTime = time.Now().String()
	CommitId = "none"

	appInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dockerhub_limit_exporter_build_info",
			Help: "Docker hub limit exporter build information.",

		},[]string{"build_time", "build_version", "commit_id"})
)

func main() {
    flag.Parse()
	appInfo.WithLabelValues(BuildTime, BuildVersion, CommitId).Add(1)
	prometheus.MustRegister(appInfo)
	prometheus.MustRegister(collector.NewHubCollector())

	http.Handle(*c.MetricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Docker hub limit exporter</title></head>
			<body>
			<h1>Docker hub limit exporter</h1>
			<p><a href="` + *c.MetricsPath + `">Metrics</a></p>
			</body>
			</html>`))

	})

	srv := &http.Server{
		Addr: *c.ListenAddr,
	}

	log.Infof("Starting docker hub limit exporter on %s", *c.ListenAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Errorf("Could not start exporter: %s", err)
		os.Exit(1)
	}
}
