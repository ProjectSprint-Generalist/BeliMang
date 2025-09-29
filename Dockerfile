########################
# Build stage
########################
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Install build deps
RUN apk add --no-cache ca-certificates git build-base

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/belimang ./

########################
# Runtime stage
########################
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /

COPY --from=builder /bin/belimang /belimang
COPY --from=builder /app/migrations /migrations

ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/belimang"]
