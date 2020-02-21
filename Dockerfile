FROM alpine:3.10.3 AS certs
RUN apk --no-cache add ca-certificates

FROM golang:1.13.5 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /nemesis .

FROM scratch
COPY --from=builder /nemesis ./
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["./nemesis"]