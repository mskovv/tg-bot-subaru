FROM golang:1.23.2

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY . .
RUN go mod tidy

ARG APP_PORT
ENV APP_PORT=${DOCKER_APP_PORT}

EXPOSE ${APP_PORT}