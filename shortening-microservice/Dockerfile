FROM golang:alpine AS builder

WORKDIR /shortening

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o shorteningserver ./cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /shortening/.env .

COPY --from=builder /shortening/shorteningserver .

CMD ["./shorteningserver"]