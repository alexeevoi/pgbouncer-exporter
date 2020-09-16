package collector

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
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
func pgbNewConnection() (*sql.DB, error) {
	psqlInfo := "postgres://pgbouncer@localhost:6432?sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func Test_NewPgbouncerCollector(t *testing.T) {
	a := assert.New(t)
	configZap()
	t.Run("success collector register", func(t *testing.T) {
		db, err := pgbNewConnection()
		s := NewPgbouncerCollector(logger, db, "pgb1")
		reg := prometheus.NewPedanticRegistry()
		if err := reg.Register(s); err != nil {
			t.Error("registration failed:", err)
		}

		a.NoError(err)
	})

}
func Test_Collect(t *testing.T) {
	a := assert.New(t)
	configZap()
	t.Run("success collecting metrics", func(t *testing.T) {
		db, _ := pgbNewConnection()
		s := NewPgbouncerCollector(logger, db, "pgb1")
		ch := make(chan prometheus.Metric)
		go func() {
			for {
				<-ch
			}
		}()
		err := s.CollectPools(ch)
		a.NoError(err)
		err = s.CollectStats(ch)
		a.NoError(err)
		err = s.CollectLists(ch)
		a.NoError(err)
	})

	t.Run("error collecting metrics", func(t *testing.T) {
		//bad connection
		db, _ := sql.Open("postgres", "")
		s := NewPgbouncerCollector(logger, db, "pgb1")
		ch := make(chan prometheus.Metric)
		go func() {
			for {
				<-ch
			}
		}()
		err := s.CollectPools(ch)
		a.Error(err)
		err = s.CollectStats(ch)
		a.Error(err)
		err = s.CollectLists(ch)
		a.Error(err)
	})

}
