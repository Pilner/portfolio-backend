FROM golang:1.26.1 AS builder

ENV GOOS=linux 
ENV CGO_ENABLED=0

WORKDIR /myapp

COPY . .

RUN go build -o myapp ./cmd/app

# Start a new minimal image
FROM alpine:3.23 as production

WORKDIR /app

COPY --from=builder /myapp/myapp .

EXPOSE 8000

CMD ["./myapp"]

