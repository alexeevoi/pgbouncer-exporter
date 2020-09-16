package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	pgbouncer_exporter "github.com/alexeevoi/pgbouncer-exporter/internal/pgbouncer-exporter"
)

var logger *zap.Logger

func configZap() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapLogger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger = zapLogger
}
func configViper() {
	viper.SetConfigName("pgbouncer-exporter")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error("Config file not found")
			panic("Config file not found")
		} else {
			logger.Sugar().Errorf("Error occurred while reading config:%v", err)
			panic("Config file error")
		}
	}
}
func main() {
	configZap()
	configViper()
	bouncersConfig := viper.Get("pgbouncers").(map[string]interface{})
	metricsHost := viper.Get("metrics_host").(string)

	pgbExporter := pgbouncer_exporter.NewExporter(logger)
	defer pgbExporter.CloseAllDb()
	errors := pgbExporter.RegisterPgbExporters(bouncersConfig)
	for _, err := range errors {
		if err != nil {
			logger.Error("Unable to register exporters", zap.Error(err))
		}
	}
	err := pgbExporter.StartExporter(metricsHost)
	if err != nil {
		logger.Error("Unable to start exporter.", zap.Error(err))
	}
}
