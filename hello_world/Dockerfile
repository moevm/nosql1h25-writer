# SOURCE CODE AND DEPENDENCIES
FROM golang:1.24.0-alpine AS builder
WORKDIR /main
COPY . .
RUN go mod download
RUN go build -o bin/main cmd/main/main.go

# FINAL STAGE
FROM alpine AS final
WORKDIR /main

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /main
USER appuser

COPY --from=builder /main/bin/main /main/main
CMD ["./main"]
