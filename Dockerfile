# Build stage

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/api/main.go

# Final stage

FROM scratch

COPY --from=builder /app/server /server

ENTRYPOINT ["./server"]
