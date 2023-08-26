# syntax=docker/dockerfile:1

FROM golang:1.21-bullseye

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /expense_bot

EXPOSE 8080

CMD [ "/expense_bot" ]