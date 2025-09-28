FROM golang:1.25 AS build-go

ARG appRelease=0.0.1

ENV APP_RELEASE=$appRelease

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w -X main.version=$APP_RELEASE" -mod=readonly -v -o app cmd/app/main.go

FROM debian:bookworm-slim

WORKDIR /app

# Install CA certificates and SQLite runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build-go /app/app app

EXPOSE 8080

CMD ["./app"]
