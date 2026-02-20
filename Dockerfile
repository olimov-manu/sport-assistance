ARG BUILDPLATFORM
ARG TARGETPLATFORM

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG CGO_ENABLED=0

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

RUN set -eux; \
    export GOOS="$TARGETOS" GOARCH="$TARGETARCH" CGO_ENABLED="$CGO_ENABLED"; \
    if [ "$TARGETARCH" = "arm" ] && [ -n "$TARGETVARIANT" ]; then export GOARM="${TARGETVARIANT#v}"; fi; \
    mkdir -p /app/bin; \
    go build -o /app/bin/sport-assistance ./cmd/main.go

FROM --platform=$TARGETPLATFORM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/sport-assistance /app/sport-assistance
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations /app/migrations

ENV GOOSE_DRIVER=postgres
ENV GOOSE_MIGRATION_DIR=/app/migrations

EXPOSE 8080

CMD ["sh", "-ec", "export GOOSE_DBSTRING=\"postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE\"; goose -dir \"$GOOSE_MIGRATION_DIR\" up && /app/sport-assistance"]
