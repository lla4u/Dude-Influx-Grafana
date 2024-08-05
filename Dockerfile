FROM golang:1.21-alpine AS builder
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/dude-cli main.go

FROM alpine:3.20
COPY --chown=65534:65534 --from=builder /go/bin/dude-cli .
COPY --chown=65534:65534 --from=builder /go/src/app/config.env .

USER 65534

ENTRYPOINT [ "/bin/ash" ]