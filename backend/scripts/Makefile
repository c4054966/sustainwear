build:
	go build -o bin/sustainwear cmd/api/main.go

run:
	go run cmd/main.go

init-db:
	./scripts/init_db.sh

run-db:
	./scripts/run.sh

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy