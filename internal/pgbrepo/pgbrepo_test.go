package pgbrepo

import (
	"database/sql"
	"math"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_Utils(t *testing.T) {
	a := assert.New(t)

	t.Run("success cast", func(t *testing.T) {

		fb := []byte("1.01")
		r, err := castDbFloat64(fb)
		a.NoError(err)
		a.Equal(r, 1.01)

		var f2 int64 = 1
		r, err = castDbFloat64(f2)
		a.NoError(err)
		a.Equal(r, float64(1))

		f3 := time.Now()
		r, err = castDbFloat64(f3)
		a.NoError(err)
		a.Equal(r, float64(f3.Unix()))

		r, err = castDbFloat64(nil)
		a.NoError(err)
		assert.True(t, math.IsNaN(r))
	})
	t.Run("error cast", func(t *testing.T) {
		s := "string"
		r, err := castDbFloat64(s)
		a.Error(err)
		assert.True(t, math.IsNaN(r))

		s2 := struct {
		}{}
		r, err = castDbFloat64(s2)
		a.Error(err)
		assert.True(t, math.IsNaN(r))

	})

	t.Run("pool_mode to index", func(t *testing.T) {
		s := "session"
		idx := poolModeIndex(s)
		a.Equal(idx, float64(1))

		s = "transaction"
		idx = poolModeIndex(s)
		a.Equal(idx, float64(2))

		s = "statement"
		idx = poolModeIndex(s)
		a.Equal(idx, float64(3))

		s = "unknown"
		idx = poolModeIndex(s)
		a.Equal(idx, float64(0))

	})

}
func pgbNewConnection() (*sql.DB, error) {
	psqlInfo := "postgres://pgbouncer@localhost:6432?sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func Test_Metrics(t *testing.T) {
	a := assert.New(t)
	db, _ := pgbNewConnection()
	t.Run("get metrics error", func(t *testing.T) {
		_, err := GetMetrics(db, "UNKNOWNCOMMAND")
		a.Error(err)
	})
	t.Run("get metrics success", func(t *testing.T) {
		metrics, err := GetMetrics(db, "POOLS")
		a.NoError(err)
		a.Equal(metrics["pgbouncerpgbouncer"].Metrics[2].Name, "cl_active")

		metrics, err = GetMetrics(db, "STATS")
		a.NoError(err)
		a.Equal(metrics["pgbouncer"].Metrics[1].Name, "total_xact_count")

		metrics, err = GetMetrics(db, "LISTS")
		a.NoError(err)
		a.Equal(metrics["LISTS"].Metrics[1].Name, "databases")
	})

}
