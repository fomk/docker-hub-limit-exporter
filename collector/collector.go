package collector

import (
	"github.com/fomk/docker-hub-limit-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	HubLimit = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dockerhub_limit_max_requests_total",
			Help: "Docker hub rate limit maximum requests.",
		})

	HubRemaining = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dockerhub_limit_remaining_requests_total",
			Help: "Docker hub rate limit remaining requests.",
		})

	upMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dockerhub_limit_up",
			Help: "Status of the last metric scrape.",
		})
)

type HubCollector struct {
	limitMetric     prometheus.Gauge
	remainingMetric prometheus.Gauge
	upMetric        prometheus.Gauge
	mutex           sync.Mutex
}


func NewHubCollector() *HubCollector  {
	return &HubCollector{
		upMetric: upMetric,
		limitMetric: HubLimit,
		remainingMetric: HubRemaining,
	}
}

func (c *HubCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upMetric.Desc()
	ch <- c.limitMetric.Desc()
	ch <- c.remainingMetric.Desc()
}

func (c *HubCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	stats, err := client.GetMetrics()
	if err != nil {
		c.upMetric.Set(0)
		ch <- c.upMetric
		log.Error(err)
		return
	}

	c.upMetric.Set(1)
	ch <- c.upMetric

	ch <- prometheus.MustNewConstMetric(c.limitMetric.Desc(),
		prometheus.GaugeValue, stats.Limit)

	ch <- prometheus.MustNewConstMetric(c.remainingMetric.Desc(),
		prometheus.GaugeValue, stats.Remaining)
}