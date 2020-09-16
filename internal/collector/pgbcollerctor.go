package collector

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/alexeevoi/pgbouncer-exporter/internal/pgbrepo"
)

const (
	namespace = "pgbouncer_exporter"
)

type pgbMetrics struct {
	Type  prometheus.ValueType
	Desc  *prometheus.Desc
	Value func(value float64) float64
}

// PgbouncerCollector is an implementation of  prometheus/client_golang/collector
// which holds main information about metrics pulled from Pgbouncer
type PgbouncerCollector struct {
	db                   *sql.DB
	logger               *zap.Logger
	pgbName              string
	errorScrapesTotal    *prometheus.Desc
	up                   prometheus.Gauge
	totalScrapes         prometheus.Counter
	totalScrapesNum      float64
	errorScrapesTotalNum float64
	metricsLists         map[string]*pgbMetrics
	metricsPools         map[string]*pgbMetrics
	metricsStats         map[string]*pgbMetrics
}

// Describe implementation of prometheus collector describe method
// Sends all possible descriptions into provided channel and returns when done
func (c *PgbouncerCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metricsLists {
		ch <- metric.Desc
	}
	for _, metric := range c.metricsPools {
		ch <- metric.Desc
	}
	for _, metric := range c.metricsStats {
		ch <- metric.Desc
	}

	ch <- c.up.Desc()
	ch <- c.totalScrapes.Desc()
}

// CollectStats collects all the metrics for PGB command SHOW STATS;
// pushes it into the channel provided by prometheus client and returns nil on success
func (c *PgbouncerCollector) CollectStats(ch chan<- prometheus.Metric) error {
	stats, err := pgbrepo.GetMetrics(c.db, "STATS")
	if err != nil {
		return err
	}
	for db, metricsList := range stats {
		for _, metricStruct := range metricsList.Metrics {
			if metric, ok := c.metricsStats[metricStruct.Name]; ok {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc,
					metric.Type,
					metric.Value(metricStruct.Value),
					db,
				)
			}
		}

	}
	return nil
}

// CollectPools collects all the metrics for PGB command SHOW POOLS;
// pushes it into the channel provided by prometheus client and returns nil on success
func (c *PgbouncerCollector) CollectPools(ch chan<- prometheus.Metric) error {
	pools, err := pgbrepo.GetMetrics(c.db, "POOLS")
	if err != nil {
		return err
	}
	for _, metricsList := range pools {
		for _, metricStruct := range metricsList.Metrics {
			if metric, ok := c.metricsPools[metricStruct.Name]; ok {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc,
					metric.Type,
					metric.Value(metricStruct.Value),
					metricStruct.Labels["db"],
					metricStruct.Labels["user"],
				)
			}
		}

	}
	return nil
}

// CollectLists collects all the metrics for PGB command SHOW LISTS;
// pushes it into the channel provided by prometheus client and returns nil on success
func (c *PgbouncerCollector) CollectLists(ch chan<- prometheus.Metric) error {

	lists, err := pgbrepo.GetMetrics(c.db, "LISTS")
	if err != nil {
		return err
	}
	for _, metricsList := range lists {
		for _, metricStruct := range metricsList.Metrics {
			if metric, ok := c.metricsLists[metricStruct.Name]; ok {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc,
					metric.Type,
					metric.Value(metricStruct.Value),
				)
			}
		}

	}
	return nil
}

// Collect implementation of prometheus collector Collect method.
// It's called by the Prometheus registry when collecting metrics.
// The implementation sends each collected metric via the provided channel
// and returns once the last metric has been sent.
func (c *PgbouncerCollector) Collect(ch chan<- prometheus.Metric) {
	c.totalScrapesNum++

	ch <- prometheus.MustNewConstMetric(
		c.totalScrapes.Desc(),
		prometheus.CounterValue,
		c.totalScrapesNum,
	)

	err := c.CollectLists(ch)
	if err != nil {
		c.logger.Error("Error collecting metrics", zap.Error(err))
		c.errorScrapesTotalNum++
	}
	err = c.CollectPools(ch)
	if err != nil {
		c.logger.Error("Error collecting metrics", zap.Error(err))
		c.errorScrapesTotalNum++
	}
	err = c.CollectStats(ch)
	if err != nil {
		c.logger.Error("Error collecting metrics", zap.Error(err))
		c.errorScrapesTotalNum++
	}

	ch <- prometheus.MustNewConstMetric(
		c.errorScrapesTotal,
		prometheus.CounterValue,
		c.errorScrapesTotalNum,
	)
}

// NewPgbouncerCollector configures a new instance of PgbouncerCollector with metrics, logger and db connection
func NewPgbouncerCollector(logger *zap.Logger, db *sql.DB, bouncerName string) *PgbouncerCollector {
	subsystem := bouncerName + "_main"
	return &PgbouncerCollector{
		db:      db,
		logger:  logger,
		pgbName: bouncerName,
		errorScrapesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_error_scrapes"),
			"Number of scrapes ended with error.",
			nil, nil,
		),
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "up"),
			Help: "Was the last scrape  successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "total_scrapes"),
			Help: "Current total scrapes.",
		}),
		metricsLists: getMetricsListsMap(bouncerName),
		metricsStats: getMetricsStatsMap(bouncerName),
		metricsPools: getMetricsPoolsMap(bouncerName),
	}
}
