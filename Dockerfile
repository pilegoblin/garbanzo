FROM golang:1.24.4-bookworm

WORKDIR /app

COPY . .

ENV PORT=8080
ENV HOST=garbanzo.chat
ENV ENVIRONMENT=prod
EXPOSE ${PORT}

RUN mkdir -p ./bin
RUN make update-tailwind
RUN make css

RUN go build -o ./bin/garbanzo ./cmd/garbanzo/*.go

CMD [ "./bin/garbanzo" ]