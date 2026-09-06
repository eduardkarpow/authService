FROM golang:1.27-alpine AS base
WORKDIR /app

RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev

# Install Air for hot-reloading
RUN go install github.com/air-verse/air@latest

COPY . .

# Expose app port
EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

FROM base AS builder

COPY . .

# Build a static binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server .

FROM gcr.io/distroless/static-debian12:nonroot AS prod

WORKDIR /app

COPY --from=builder /app/server /app/server

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/server"]