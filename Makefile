.PHONY: build
build: gen css
	@go build -o ./bin/garbanzo ./cmd/garbanzo/*.go

.PHONY: run
run: build
	@./bin/garbanzo

.PHONY: gen
gen:
	@rm -rf ./db/sqlc && sqlc generate

.PHONY: update-tailwind
update-tailwind:
	@curl -sL0 https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 > ./bin/tailwind
	@chmod +x ./bin/tailwind

.PHONY: css
css:
	@./bin/tailwind -i ./templates/main.css -o ./public/output.css --minify