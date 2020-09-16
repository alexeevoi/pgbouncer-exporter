package pgbouncer_exporter

import (
	"database/sql"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/alexeevoi/pgbouncer-exporter/internal/collector"
)

type Exporter struct {
	dbs    map[string]*sql.DB
	logger *zap.Logger
}

// Opens new Pgbouncer connection
func (e *Exporter) pgbNewConnection(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Returns new instance of pgbouncer exporter
func NewExporter(logger *zap.Logger) *Exporter {
	return &Exporter{dbs: make(map[string]*sql.DB), logger: logger}
}

// Closes all pgbouncer connections
func (e *Exporter) CloseAllDb() []error {
	var errors []error
	for _, db := range e.dbs {
		errors = append(errors, db.Close())
	}
	return errors
}

// RegisterPgbExporters opens new connection to pgbouncer for each entity described in config map
// Creates new pgbouncer collector instance and registers it in prometheus client.
func (e *Exporter) RegisterPgbExporters(bouncersConfig map[string]interface{}) []error {
	var errors []error
	for bouncerName, bouncerConfig := range bouncersConfig {
		config := bouncerConfig.(map[string]interface{})
		db, err := e.pgbNewConnection(config["connection"].(string))
		if err != nil {
			errors = append(errors, err)
			e.logger.Error("Unable to connect ", zap.Error(err))
		}

		collector := collector.NewPgbouncerCollector(e.logger, db, bouncerName)
		prometheus.MustRegister(collector)
		e.dbs[bouncerName] = db

	}
	return errors
}

// StartExporter creates http handle for prometheus on provided host ex. localhost:8080
func (e *Exporter) StartExporter(metricsHost string) error {
	http.Handle("/metrics", promhttp.Handler())
	e.logger.Sugar().Infof("Beginning to serve on %v", metricsHost)
	err := http.ListenAndServe(metricsHost, nil)
	if err != nil {
		return err
	}
	return nil
}
