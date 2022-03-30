FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download
RUN go mod verify

FROM build AS rest
COPY . ./
RUN go build -v -o server ./cmd/rest

FROM debian:buster-slim as final-rest
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=rest /app/server /app/server

EXPOSE 1234

CMD ["/app/server"]