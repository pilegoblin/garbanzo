FROM golang:latest

WORKDIR /app

COPY . .

ENV PORT=8080
ENV HOST=garbanzo.chat
ENV ENVIRONMENT=prod
EXPOSE ${PORT}

RUN make update-tailwind
RUN make css

RUN go build -o ./bin/garbanzo ./cmd/garbanzo/*.go

CMD [ "./bin/garbanzo" ]