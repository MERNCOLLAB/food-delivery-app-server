FROM golang:1.24.5-alpine AS builder

WORKDIR /delivery-app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /delivery-app

COPY --from=builder /delivery-app/server .

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./server"]