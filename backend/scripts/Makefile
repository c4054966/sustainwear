build:
	go build -o sustainwear.exe

run:
	go run main.go

init-db:
	./scripts/init_db.sh

run-db:
	./scripts/run.sh

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy