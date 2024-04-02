FROM golang:1.20.14 as builder-plugin

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -buildmode=plugin -o krakend_revoke_jwt.so ./cmd/krakend-jwt-revoker/main.go

FROM devopsfaith/krakend:latest
COPY --from=builder-plugin /app/krakend_revoke_jwt.so /etc/krakend/plugins/krakend_revoke_jwt.so