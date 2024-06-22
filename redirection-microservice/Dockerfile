FROM golang:alpine AS builder

WORKDIR /redirection

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o redirectionserver ./cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /redirection/.env .

COPY --from=builder /redirection/redirectionserver .

CMD ["./redirectionserver"]