all: cmd

dep:
	rm -rf vendor/
	dep ensure -v

.PHONY: cmd
cmd:
	go build cmd/hakagi/hakagi.go
