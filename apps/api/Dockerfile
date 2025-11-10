# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/bookshop ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates curl
ENV PORT=8080
ENV DB_DSN=postgres://postgres:postgres@postgres:5432/books?sslmode=disable
COPY --from=builder /bin/bookshop /bin/bookshop
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD curl -f http://localhost:8080/books || exit 1
ENTRYPOINT ["/bin/bookshop"]
