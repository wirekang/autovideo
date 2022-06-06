dev:
	go run main.go $(ARGS)

test:
	go test -v ./...

sample:
	cd sample; go run ../main.go --config=config.json --script=script.txt --audios=audios --output=output.mp4

sample-debug:
	cd sample; go run ../main.go --config=config.json --script=script.txt --audios=audios --output=output.mp4 --debug

.PHONY: sample test dev