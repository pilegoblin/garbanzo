build-auth: db
	@go build -o ./bin/auth ./cmd/auth/*.go

run-auth: build-auth
	@./bin/auth

db:
	@sqlc generate