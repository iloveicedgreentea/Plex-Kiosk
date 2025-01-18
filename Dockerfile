# Frontend build stage
FROM oven/bun:1 as frontend-builder
WORKDIR /frontend
COPY frontend/package.json frontend/bun.lockb ./
RUN bun install
COPY frontend/ .
RUN bun run build

# Backend build stage
FROM golang:alpine as backend-builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd .
RUN GOOS=linux go build -o plex-viewer .

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

# Copy the frontend build
COPY --from=frontend-builder /frontend/dist /usr/share/nginx/html

# Copy compiled backend
COPY --from=backend-builder /app/plex-viewer .
RUN chown appuser:appgroup plex-viewer && chmod +x plex-viewer


# Copy startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 80
ENV GIN_MODE=release
ENV ALLOWED_LIBRARIES="Movies,TV Shows"

ENTRYPOINT ["/start.sh"]