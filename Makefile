.PHONY: build
build: gen
	@go build -o ./bin/garbanzo ./cmd/garbanzo/*.go

.PHONY: run
run: build gen
	@./bin/garbanzo

.PHONY: gen
gen:
	rm -rf ./db/sqlc && sqlc generate
