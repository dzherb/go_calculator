GOBIN ?= $$(go env GOPATH)/bin

.PHONY: test
test:
	@go test ./... -count=1

.PHONY: cover
cover:
	@go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

.PHONU: clean
clean:
	@rm -f cover.out

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: fmt
fmt:
	@golangci-lint fmt ./...