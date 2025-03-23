.PHONY: build
build: gen
	@go build -o ./bin/auth ./cmd/auth/*.go

.PHONY: run
run: build
	@./bin/auth

.PHONY: gen
generate: db
	@go generate ./ent
