FROM golang:1.23.2 AS builder

WORKDIR /github.com/mskovv/tg-bot-subaru96

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache
COPY ./go.* ./

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go install github.com/air-verse/air@latest

RUN go build -o app ./cmd/main.go
FROM golang:1.23.2

RUN apt-get update && apt-get install -y bash && apt-get clean

WORKDIR /usr/src/app

COPY --from=builder /github.com/mskovv/tg-bot-subaru96 .
COPY --from=builder /go/bin/air /usr/local/bin/air

RUN which air && chmod +x /usr/local/bin/air

EXPOSE 9000