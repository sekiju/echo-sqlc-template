FROM golang:1.21-alpine3.19 AS build
WORKDIR /go/src/echo-sqlc-template

ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

COPY .. .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/echo-sqlc-template cmd/api/main.go

FROM alpine:3.19
COPY --from=build /go/bin/echo-sqlc-template /usr/bin/local/echo-sqlc-template
COPY --from=build /go/src/echo-sqlc-template/resources/migrations /usr/bin/local/migrations

ENV DATABASE_MIGRATIONS /usr/bin/local/migrations

EXPOSE 8000
ENTRYPOINT ["/usr/bin/local/echo-sqlc-template"]