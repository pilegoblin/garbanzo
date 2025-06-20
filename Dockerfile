FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN go build -o ./bin/garbanzo ./cmd/garbanzo/*.go

ENV PORT=8080
ENV HOST=garbanzo.chat
ENV ENVIRONMENT=prod
EXPOSE ${PORT}

CMD [ "./bin/garbanzo" ]