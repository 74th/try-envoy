FROM golang:1.14 AS build-env
WORKDIR /go/src/github.com/74th/try-envoy/router
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
WORKDIR /go/src/github.com/74th/try-envoy/router/server
RUN go build -o server server.go
RUN go run /usr/local/go/src/crypto/tls/generate_cert.go -duration 876000h0m0s -host envoy

FROM debian
WORKDIR /app
COPY --from=build-env \
 /go/src/github.com/74th/try-envoy/router/server/server \
 /go/src/github.com/74th/try-envoy/router/server/cert.pem \
 /go/src/github.com/74th/try-envoy/router/server/key.pem \
 ./
EXPOSE 443
ENTRYPOINT ./server
