FROM golang:1.23-alpine as builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY ./cmd .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o plex-viewer .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install nginx and dependencies
RUN apk add --no-cache nginx curl gettext

# Create appuser and group
RUN addgroup -S appgroup && \
    adduser -S -G appgroup -s /bin/bash appuser && \
    mkdir -p /var/cache/nginx /var/log/nginx /run/nginx /app/cache && \
    chown -R nginx:nginx /var/cache/nginx /var/log/nginx /run/nginx && \
    chown -R appuser:appgroup /app/cache

# Copy nginx configuration
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.template

# Copy compiled application
COPY --from=builder /app/plex-viewer .
RUN chown appuser:appgroup plex-viewer && chmod +x plex-viewer

# Copy templates
COPY templates/ templates/
RUN chown -R appuser:appgroup templates

# Copy startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 80
ENV GIN_MODE=release

ENTRYPOINT ["/start.sh"]