# syntax=docker/dockerfile:1.6

FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY openapi ./openapi
COPY domain ./domain
COPY internal ./internal
COPY handlers ./handlers
COPY cmd ./cmd
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/api ./api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
