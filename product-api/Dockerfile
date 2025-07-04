# ---------- Builder stage ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o product-api main.go

# ---------- Runtime stage ----------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/product-api .
EXPOSE 8080
CMD ["./product-api"]
