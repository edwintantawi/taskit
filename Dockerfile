FROM golang:1.19-alpine3.17 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go

FROM scratch

COPY --from=builder /app/main .

CMD ["./main"]