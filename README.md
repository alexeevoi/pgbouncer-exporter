
# pgbouncer-exporter  
Экспортирует метрики pgbouncer в prometheus

Выводит все метрики по <code>SHOW LISTS/POOLS/STATS</code>  
  
Настройка коннектов и ручек в файле pgbouncer-exporter.yml  

    pgbouncers:  
      pgb1:  
        connection: "postgres://pgbouncer@localhost:6432?sslmode=disable"  
      ody1:  
        connection: "postgres://odyssey@localhost:6532?sslmode=disable"  
      pgb2:  
        connection: "postgres://pgbouncer@localhost123:6431?sslmode=disable"  
    metrics_host: "localhost:8080"

Метрики выводятся в формате: 

    pgbouncer_exporter_ALIAS_COMMAND_METRICNAME

  
Pooling_mode выводится в формате:  

    HELP pgbouncer_exporter_pgb1_pools_pool_mode The pooling mode in use. 1 - session, 2 - transaction, 3 - statement  
    TYPE pgbouncer_exporter_pgb1_pools_pool_mode gauge  
          pgbouncer_exporter_pgb1_pools_pool_mode{db="pgbouncer",user="pgbouncer"} 3   

  
Для тестов поднимаются контейнеры докера:        `make run  make test`  


  

    user@user:~/go/src/pgbouncer-exporter$ make test  
    docker-compose -f docker/docker-compose.yml down  
    Stopping docker_pgbouncer_1 ... done  
    Stopping docker_postgres_1  ... done  
    Stopping docker_odyssey_1   ... done  
    Removing docker_pgbouncer_1 ... done  
    Removing docker_postgres_1  ... done  
    Removing docker_odyssey_1   ... done  
    Removing network docker_default  
    docker-compose -f docker/docker-compose.yml up -d  
    Creating network "docker_default" with the default driver  
    Creating docker_postgres_1 ... done  
    Creating docker_odyssey_1   ... done  
    Creating docker_pgbouncer_1 ... done  
    go test ./... -coverprofile c.out  
    ?       git.ozon.dev/oalekseev/pgbouncer-exporter/cmd/pgbouncer-exporter        [no test files]  
    ok      git.ozon.dev/oalekseev/pgbouncer-exporter/internal/collector    0.037s  coverage: 83.5% of statements  
    ok      git.ozon.dev/oalekseev/pgbouncer-exporter/internal/pgbouncer-exporter   0.008s  coverage: 84.6% of statements  
    ok      git.ozon.dev/oalekseev/pgbouncer-exporter/internal/pgbrepo      0.013s  coverage: 86.5% of statements  

