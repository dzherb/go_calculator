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

.PHONY: proto
proto:
	@./scripts/gen_proto.sh

.PHONY: migrations
migrations:
	@./scripts/migrations.sh

.PHONY: migrate
migrate:
	@[ -n "$$DATABASE_URL" ] || (echo "Error: DATABASE_URL is unset" && exit 1)
	@migrate -database $DATABASE_URL -path internal/storage/migrations up

.PHONY: schema
schema:
	@./scripts/schema.sh