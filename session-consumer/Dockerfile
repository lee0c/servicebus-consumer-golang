from golang:1.12-alpine as builder

WORKDIR "/app"

RUN apk --update add git ca-certificates

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch as binary
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/session-consumer /consumer

ENTRYPOINT ["/consumer"]