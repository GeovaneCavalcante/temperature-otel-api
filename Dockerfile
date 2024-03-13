FROM golang:1.21 as build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server cmd/server.go

FROM alpine:3.14

WORKDIR /app

COPY --from=build /app/bin/server .
COPY --from=build /app/.env .
COPY --from=build /app/docker-compose-files/docker-entrypoint.sh .


ENTRYPOINT ["/bin/sh", "docker-entrypoint.sh"]