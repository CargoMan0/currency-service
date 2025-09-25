# Stage 1: Build
FROM golang:1.23.1-alpine AS builder
WORKDIR /app

COPY . .

ARG BUILD_TARGET
RUN go build -o /app/bin/app $BUILD_TARGET

# Stage 2: Create a smaller image for running the binary (Optimization)
FROM scratch

COPY --from=builder /app/bin/app /

ENTRYPOINT ["/app"]