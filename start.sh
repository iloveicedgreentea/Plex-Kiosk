#!/bin/bash
set -e

echo "Starting Plex Library Viewer..."
echo "PLEXURL is set to: $PLEXURL"

# Create required nginx directories
mkdir -p /var/log/nginx
chown -R nginx:nginx /var/log/nginx

# Replace environment variables in nginx config
echo "Configuring nginx..."
envsubst '$PLEXURL' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

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

echo "Starting FastAPI application with Gunicorn..."
# Start the FastAPI application with Gunicorn
exec gunicorn app:app \
    --workers 4 \
    --worker-class uvicorn.workers.UvicornWorker \
    --bind 127.0.0.1:8000 \
    --user appuser \
    --access-logfile - \
    --error-logfile - \
    --log-level info