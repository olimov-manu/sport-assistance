FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN for i in 1 2 3 4 5; do go mod download && exit 0; echo "go mod download failed (attempt $i), retrying..."; sleep 3; done; exit 1

COPY . .

RUN mkdir -p /app/bin && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/sport-assistance ./cmd/main.go

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/sport-assistance /app/sport-assistance

EXPOSE 8080

CMD ["/app/sport-assistance"]
