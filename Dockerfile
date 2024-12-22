FROM golang:1.23.2 AS builder

WORKDIR /usr/src/app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

FROM golang:1.23.2

WORKDIR /usr/src/app
RUN go install github.com/air-verse/air@latest


COPY --from=builder /usr/src/app/app .
COPY . .

EXPOSE 9000
CMD ["./app"]