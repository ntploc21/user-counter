FROM golang:1.25.1-alpine as builder
RUN apk add build-base

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o core cmd/*.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/core ./

CMD ["/app/core", "serve"]
