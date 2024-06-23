FROM golang:alpine AS builder

WORKDIR /gateway

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o gatewayserver ./cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /gateway/.env .

COPY --from=builder /gateway/gatewayserver .

CMD ["./gatewayserver"]