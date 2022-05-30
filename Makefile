dev:
	go run main.go $(ARGS)

test:
	go test -v tests/* && go test -v