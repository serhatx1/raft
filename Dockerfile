# syntax=docker/dockerfile:1
FROM golang:1.23.5 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o raft main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/raft ./raft
EXPOSE 8001 8002 8003
CMD ["./raft"] 