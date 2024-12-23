FROM python:3.11-slim

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# Install system dependencies and nginx
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    nginx \
    gettext-base \
    && rm -rf /var/lib/apt/lists/*

# Create users and set up permissions
RUN useradd -m -s /bin/bash appuser && \
    # Create nginx user and group if they don't exist
    adduser --system --no-create-home --shell /bin/false --group --disabled-login nginx

# Copy nginx configuration
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.template

# Create necessary directories and set permissions
RUN mkdir -p /var/cache/nginx /var/log/nginx /var/lib/nginx /run/nginx /app/cache && \
    chown -R nginx:nginx /var/cache/nginx /var/log/nginx /var/lib/nginx /run/nginx && \
    chown -R appuser:appuser /app/cache

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY . .

# Final permissions adjustments
RUN chown -R appuser:appuser /app && \
    chmod -R g+w /app/cache

# Expose port 80
EXPOSE 80

# Copy the startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Use the startup script as the entrypoint
ENTRYPOINT ["/start.sh"]