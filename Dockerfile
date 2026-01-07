# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./
COPY templates/ ./templates/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o heartflow-demo .

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/heartflow-demo .

EXPOSE 8080

ENV PORT=8080

USER nobody:nobody

ENTRYPOINT ["/app/heartflow-demo"]
