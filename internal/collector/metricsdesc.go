package collector

import "github.com/prometheus/client_golang/prometheus"

/*
This functions return maps of prometheus descriptions (required for custom exporter)
for different PgBouncer SHOW *; commands
*/
func getMetricsListsMap(bouncerName string) map[string]*pgbMetrics {
	subsystem := bouncerName + "_lists"
	return map[string]*pgbMetrics{
		"databases": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "databases"),
				"Count of databases.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"users": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "users"),
				"Count of users.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"pools": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "pools"),
				"Count of pools.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"free_clients": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "free_clients"),
				"Count of free_clients.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"used_clients": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "used_clients"),
				"Count of used clients.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"login_clients": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "login_clients"),
				"Count of clients in login state",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"free_servers": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "free_servers"),
				"Count of free_servers.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"used_servers": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "used_servers"),
				"Count of used_servers.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"dns_names": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "dns_names"),
				"Count of DNS names in the cache.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"dns_zones": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "dns_zones"),
				"Count of DNS zones in the cache.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"dns_queries": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "dns_queries"),
				"Count of in-flight DNS queries.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"dns_pending": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "dns_pending"),
				"Count of dns_pending.",
				nil, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
	}
}
func getMetricsPoolsMap(bouncerName string) map[string]*pgbMetrics {
	subsystem := bouncerName + "_pools"
	return map[string]*pgbMetrics{
		"cl_active": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "cl_active"),
				"Client connections that are linked to server connection and can process queries.",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"cl_waiting": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "cl_waiting"),
				"Client connections that have sent queries but have not yet got a server connection.",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"sv_active": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "sv_active"),
				"Server connections that are linked to a client.",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"sv_idle": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "sv_idle"),
				"Server connections that are unused and immediately usable for client queries.",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"sv_used": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "sv_used"),
				"Server connections that have been idle for more than server_check_delay",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"sv_tested": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "sv_tested"),
				"Server connections that are currently running either server_reset_query or server_check_query",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"sv_login": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "sv_login"),
				"Server connections currently in the process of logging in.",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"maxwait": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "maxwait"),
				"How long the first (oldest) client in the queue has waited, in seconds.",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"maxwait_us": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "maxwait_us"),
				"Microsecond part of the maximum waiting time.",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"pool_mode": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "pool_mode"),
				"The pooling mode in use. 1 - session, 2 - transaction, 3 - statement",
				[]string{"db", "user"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
	}
}
func getMetricsStatsMap(bouncerName string) map[string]*pgbMetrics {
	subsystem := bouncerName + "_stats"
	return map[string]*pgbMetrics{
		"total_xact_count": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "total_xact_count"),
				"Total number of SQL transactions pooled by pgbouncer.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"total_query_count": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "total_query_count"),
				"Total number of SQL queries pooled by pgbouncer.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"total_received": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "total_received"),
				"Total volume in bytes of network traffic received by pgbouncer.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"total_sent": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "total_sent"),
				"Total volume in bytes of network traffic sent by pgbouncer.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"total_xact_time": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "total_xact_time"),
				"Total number of microseconds spent by pgbouncer when connected to PostgreSQL in a transaction, either idle in transaction or executing queries.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"total_query_time": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "total_query_time"),
				"Total number of microseconds spent by pgbouncer when actively connected to PostgreSQL, executing queries.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"total_wait_time": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "total_wait_time"),
				"Time spent by clients waiting for a server, in microseconds.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"avg_xact_count": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "avg_xact_count"),
				"Average transactions per second in last stat period.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"avg_query_count": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "avg_query_count"),
				"Average queries per second in last stat period.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"avg_recv": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "avg_recv"),
				"Average received (from clients) bytes per second.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"avg_sent": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "avg_sent"),
				"Average sent (to clients) bytes per second.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"avg_xact_time": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "avg_xact_time"),
				"Average transaction duration, in microseconds.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"avg_query_time": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "avg_query_time"),
				"Average query duration, in microseconds.",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
		"avg_wait_time": {
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "avg_wait_time"),
				"Time spent by clients waiting for a server, in microseconds (average per second).",
				[]string{"db"}, nil,
			),
			Value: func(value float64) float64 {
				return value
			},
		},
	}
}
