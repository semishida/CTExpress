FROM golang:alpine AS builder

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/ctexpress ./kodmain.go

COPY ./HTML /app/HTML

FROM alpine

WORKDIR /app

COPY --from=builder /app /app

CMD ["/app/ctexpress"]
