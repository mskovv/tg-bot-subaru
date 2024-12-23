FROM golang:1.23.2 AS builder

WORKDIR /github.com/mskovv/tg-bot-subaru96

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/air-verse/air@latest

FROM golang:1.23.2

RUN apt-get update && apt-get install -y bash && apt-get clean

WORKDIR /usr/src/app

COPY --from=builder /github.com/mskovv/tg-bot-subaru96 .
COPY --from=builder /go/bin/air /usr/local/bin/air

RUN which air && chmod +x /usr/local/bin/air

EXPOSE 9000