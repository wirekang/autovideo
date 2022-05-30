dev:
	go run main.go $(ARGS)

test:
	go test -v ./...