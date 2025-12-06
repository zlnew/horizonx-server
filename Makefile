APP_NAME=horizonx-server
ENTRY=./cmd/horizonx-server

build:
	go build -o bin/$(APP_NAME) $(ENTRY)

run:
	go run $(ENTRY)

clean:
	rm -rf bin/$(APP_NAME)

