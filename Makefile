all: cmd

dep:
	rm -rf vendor/
	dep ensure -v

cmd:
	go build cmd/hakagi/hakagi.go
