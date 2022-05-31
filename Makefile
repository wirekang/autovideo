dev:
	go run main.go $(ARGS)

test:
	go test -v tests/* && go test -v

sample:
	cd sample; go run ../main.go --config=config.json --txt=txt.txt --audios=audios --output=output.mp4

.PHONY: sample test dev