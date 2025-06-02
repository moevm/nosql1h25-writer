# Step 1: Modules caching
FROM golang:1.24-alpine3.22 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.24-alpine3.22 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

# Step 3: Final
FROM alpine:3.22
RUN apk add --no-cache \
    mongodb-tools=100.9.4-r10 \
    gzip=1.14-r1 \
    bash=5.2.37-r0 \
    ca-certificates=20241121-r2
COPY dump.gzip /dump.gzip
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
CMD ["/app"]
