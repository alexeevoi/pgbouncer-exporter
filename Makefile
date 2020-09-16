build:
	mkdir -p bin
	go build -o bin/pgbouncer-exporter cmd/pgbouncer-exporter/main.go

run:
	docker-compose -f docker/docker-compose.yml down
	docker-compose -f docker/docker-compose.yml up -d
	go run cmd/pgbouncer-exporter/main.go

test:
	docker-compose -f docker/docker-compose.yml down
	docker-compose -f docker/docker-compose.yml up -d
	go test -race ./... -coverprofile c.out
docker_down:
	docker-compose -f docker/docker-compose.yml down
