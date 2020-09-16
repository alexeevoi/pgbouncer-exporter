package pgbouncer_exporter

import (
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func configZap() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zaplogger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger = zaplogger
}
func configViper(filename string) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error("Config file not found")
		} else {
			logger.Sugar().Errorf("Error occurred while reading config:%v", err)
		}
	}
}

func Test_Exporter(t *testing.T) {
	a := assert.New(t)
	configZap()

	t.Run("success exporter creation", func(t *testing.T) {
		configViper("pgbouncer-exporter-test")
		bouncersConfig := viper.Get("pgbouncers").(map[string]interface{})
		exp := NewExporter(logger)
		a.NotEmpty(exp)

		errors := exp.RegisterPgbExporters(bouncersConfig)
		for _, err := range errors {
			a.NoError(err)
		}
		errors = exp.CloseAllDb()
		for _, err := range errors {
			a.NoError(err)
		}

	})
	t.Run("failure exporter start", func(t *testing.T) {
		configViper("pgbouncer-exporter-test-bad")
		metricsHost := viper.Get("metrics_host").(string)
		exp := NewExporter(logger)
		a.NotEmpty(exp)

		err := exp.StartExporter(metricsHost)
		a.Error(err)
	})

}
