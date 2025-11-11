FROM golang:1.24 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api

FROM gcr.io/distroless/base-debian12 AS runtime

ENV PORT=8080
EXPOSE 8080

COPY --from=build /bin/api /bin/api

ENTRYPOINT ["/bin/api"]

