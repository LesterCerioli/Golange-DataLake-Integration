# Stage 1: Build
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o dataa-lake .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dataa-lake .

CMD ["./dataa-lake"]

EXPOSE 3000
