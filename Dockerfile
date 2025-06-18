FROM golang:alpine3.21 AS build

WORKDIR /build
ENV PATH="/usr/local/go/bin:$PATH"
ENV CGO_ENABLED=0

# Install dependencies
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev \
    ca-certificates \
    bash \
    git
RUN mkdir -p /etc/ssl/certs && \
    update-ca-certificates --fresh

# Install backend dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . ./

# Build 
RUN go build -o /build/main ./cmd/main

# Step 3: Final Image
FROM alpine:3.21
# FROM scratch

WORKDIR /bot

# Copy built frontend and backend
COPY --from=build /build/main /app/main

EXPOSE 3000
CMD ["/app/main"]
