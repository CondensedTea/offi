FROM golang:1.17-alpine as builder

WORKDIR /build

RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 go build -o bin/offi ./cmd/api/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bin/offi .

ENTRYPOINT ["./offi"]