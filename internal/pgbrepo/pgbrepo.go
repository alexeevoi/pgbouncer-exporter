package pgbrepo

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

type MetricStruct struct {
	Name    string
	IsLabel bool
	Value   float64
	Labels  map[string]string
}

type MetricsList struct {
	Name    string
	Metrics map[int]MetricStruct
}

type MetricsMap map[string]*MetricsList

func (m *MetricsMap) Add(name string) {
	(*m)[name] = &MetricsList{Name: name, Metrics: make(map[int]MetricStruct)}
}

// Casts value returned by pgbouncer to float64 for using in prometheus
func castDbFloat64(t interface{}) (float64, error) {
	switch v := t.(type) {
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case time.Time:
		return float64(v.Unix()), nil
	case []byte:
		strV := string(v)
		result, err := strconv.ParseFloat(strV, 64)
		if err != nil {
			return math.NaN(), fmt.Errorf(" %v ", err)
		}
		return result, nil
	case nil:
		return math.NaN(), nil
	default:
		return math.NaN(), fmt.Errorf("unable to cast %v", v)
	}
}
func poolModeIndex(poolMode string) float64 {
	switch poolMode {
	case "session":
		return 1
	case "transaction":
		return 2
	case "statement":
		return 3
	default:
		return 0
	}
}

// Returns map of metrics for certain pgbouncer SHOW *; command
func GetMetrics(db *sql.DB, command string) (map[string]*MetricsList, error) {
	switch command {
	case "LISTS":
		return GetMetricsByRows(db, command)
	case "POOLS":
		return GetMetricsPools(db, command)
	case "STATS":
		return GetMetricsByCols(db, command)
	}
	return nil, fmt.Errorf("unknown command %v", command)
}

// Executes pgbouncer SHOW *; command and returns sql.Rows, column names , data interface to store results
// and scanPointers to use with sql.rows.scan(&data) OR sql error
func execQuery(db *sql.DB, command string) (*sql.Rows, []string, []interface{}, []interface{}, error) {
	sqlStatement := fmt.Sprintf(`SHOW %v;`, command)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var data = make([]interface{}, len(columns))
	var scanPointers = make([]interface{}, len(columns))
	for i := range data {
		scanPointers[i] = &data[i]
	}
	return rows, columns, scanPointers, data, nil
}

// Returns map of metrics for pgb collector for command SHOW POOLS: OR error
// Unique metrics key is combination of username and database name
// pool_mode(string) is casted to float64 with poolModeIndex function
func GetMetricsPools(db *sql.DB, command string) (map[string]*MetricsList, error) {
	metricsList := make(MetricsMap)
	rows, columns, scanPointers, data, err := execQuery(db, command)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dbname string
		var userlabel string
		var metricKey string
		err := rows.Scan(scanPointers...)
		if err != nil {
			return nil, err
		}

		for idx, column := range columns {
			if column == "database" {
				dbname = data[idx].(string)
			}
			if column == "user" {
				userlabel = data[idx].(string)
				metricKey = dbname + userlabel
				metricsList.Add(metricKey)

			}
			if column == "pool_mode" {
				poolMode := data[idx].(string)
				poolModeIdx := poolModeIndex(poolMode)
				metricsList[metricKey].Metrics[idx] = MetricStruct{Name: column, IsLabel: true, Value: poolModeIdx, Labels: map[string]string{"user": userlabel, "db": dbname}}

			}
			if reflect.TypeOf(data[idx]).String() != "string" {
				value, _ := castDbFloat64(data[idx])
				metricsList[metricKey].Metrics[idx] = MetricStruct{Name: column, Value: value, Labels: map[string]string{"user": userlabel, "db": dbname}}
			}
		}
	}
	if len(metricsList) == 0 {
		return nil, fmt.Errorf("pgb returned zero rows")
	}

	return metricsList, nil
}

// Returns map of metrics for pgb collector for command SHOW STATS: OR error
// Unique metrics key is  database name
// All the metrics are labeled with "db":database name
func GetMetricsByCols(db *sql.DB, command string) (map[string]*MetricsList, error) {
	metricsList := make(MetricsMap)
	rows, columns, scanPointers, data, err := execQuery(db, command)

	if err != nil {
		return nil, err
	}
	var dbname string
	for rows.Next() {
		err := rows.Scan(scanPointers...)
		if err != nil {
			return nil, err
		}

		for idx, column := range columns {
			if column == "database" {
				dbname = data[idx].(string)
				metricsList.Add(dbname)
			}
			if reflect.TypeOf(data[idx]).String() != "string" {
				value, _ := castDbFloat64(data[idx])
				metricsList[dbname].Metrics[idx] = MetricStruct{Name: column, Value: value, Labels: map[string]string{"db": dbname}}
			}
		}
	}
	if len(metricsList) == 0 {
		return nil, fmt.Errorf("pgb returned zero rows")
	}

	return metricsList, nil
}

// Returns map of metrics for pgb collector for command SHOW LISTS: OR error
// Unique metrics key is command name
func GetMetricsByRows(db *sql.DB, command string) (map[string]*MetricsList, error) {

	metricsList := make(MetricsMap)
	metricsList.Add(command)
	rows, columns, scanPointers, data, err := execQuery(db, command)

	if err != nil {
		return nil, err
	}

	for rowIdx := 0; rows.Next(); {
		rowIdx++
		err := rows.Scan(scanPointers...)
		if err != nil {
			return nil, err
		}
		var name string
		var value float64
		for idx := range columns {

			if reflect.TypeOf(data[idx]).String() == "string" {
				name = data[idx].(string)
			} else {
				value, err = castDbFloat64(data[idx])
			}
			if err != nil {
				return nil, err
			}

		}
		metricsList[command].Metrics[rowIdx] = MetricStruct{Name: name, Value: value}
	}

	if len(metricsList) == 0 {
		return nil, fmt.Errorf("pgb returned zero rows")
	}
	return metricsList, nil
}
