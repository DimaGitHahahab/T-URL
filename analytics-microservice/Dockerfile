FROM golang:alpine AS builder

WORKDIR /analytics

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o analyticsserver ./cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /analytics/.env .

COPY --from=builder /analytics/analyticsserver .

CMD ["./analyticsserver"]