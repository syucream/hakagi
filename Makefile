all: cmd

.PHONY: dep
dep:
	rm -rf vendor/
	dep ensure -v

.PHONY: cmd
cmd:
	go build cmd/hakagi/hakagi.go

.PHONY: test
test:
	go test github.com/syucream/hakagi/...

# Set GITHUB_TOKEN personal access token and create release git tag
.PHONY: release
release:
	go get -u github.com/goreleaser/goreleaser
	goreleaser --rm-dist
