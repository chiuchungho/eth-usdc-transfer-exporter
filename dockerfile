FROM golang:1.22 as build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOSUMDB=off

WORKDIR /src

COPY . .

# build service
RUN go build -mod=vendor -a -tags netgo --installsuffix netgo -ldflags '-w' -o eth-usdc-transfer-exporter cmd/main.go

# build run container
FROM alpine:3.14
RUN apk --no-cache add \
    ca-certificates

RUN adduser -D -s /bin/sh app

COPY --from=build /src/eth-usdc-transfer-exporter /bin/eth-usdc-transfer-exporter
RUN chmod a+x /bin/eth-usdc-transfer-exporter

USER app