#!/bin/sh
set -e

echo "Starting Plex Library Viewer..."
echo "PLEX_URL is set to: $PLEX_URL"

# Create required nginx directories
mkdir -p /var/log/nginx
chown -R nginx:nginx /var/log/nginx

# Replace environment variables in nginx config
echo "Configuring nginx..."
envsubst '$PLEX_URL' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

# Test nginx config
echo "Testing nginx configuration..."
nginx -t

# Start Nginx in the background
echo "Starting nginx..."
nginx -g 'daemon off;' &
NGINX_PID=$!

# Wait a moment for nginx to start
sleep 2

# Check if nginx is running
if ! kill -0 $NGINX_PID 2>/dev/null; then
    echo "Nginx failed to start - check nginx error log"
    cat /var/log/nginx/error.log
    exit 1
fi

echo "Starting Go application..."
# Start the Go application
exec ./plex-viewer