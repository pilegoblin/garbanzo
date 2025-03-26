.PHONY: build
build: gen
	@go build -o ./bin/auth ./cmd/auth/*.go

.PHONY: run
run: build gen
	@./bin/auth

.PHONY: gen
gen:
	@go generate ./ent
