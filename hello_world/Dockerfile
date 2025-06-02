# SOURCE CODE AND DEPENDENCIES
FROM golang:1.24-alpine3.22 AS builder
WORKDIR /main
COPY . .
RUN go mod download && go build -o bin/main cmd/main/main.go

# FINAL STAGE
FROM alpine:3.22 AS final
WORKDIR /main

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && chown -R appuser:appgroup /main
USER appuser

COPY --from=builder /main/bin/main /main/main
CMD ["./main"]
