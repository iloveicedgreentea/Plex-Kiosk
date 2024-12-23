import os
import asyncio
import logging
from datetime import datetime
from typing import Dict
import json
from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.templating import Jinja2Templates
from plexapi.server import PlexServer
from aiocache import caches
from pathlib import Path

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(title="Plex Library Viewer")

# Setup templates
templates = Jinja2Templates(directory="templates")

# Configure cache
caches.set_config(
    {
        "default": {
            "cache": "aiocache.SimpleMemoryCache",
            "serializer": {"class": "aiocache.serializers.JsonSerializer"},
        }
    }
)
cache = caches.get("default")

# Environment variables
PLEX_URL = os.getenv("PLEXURL", "http://localhost:32400")
ALLOWED_LIBRARIES = os.getenv("ALLOWED_LIBRARIES", "").split(",")
REFRESH_INTERVAL = int(os.getenv("REFRESH_INTERVAL", "21600"))  # 6 hours

# Cache file path
CACHE_FILE = Path("/app/cache/library_data.json")


async def fetch_library_data() -> Dict:
    """Fetch library data from Plex server"""
    try:
        logger.info("Fetching library data from Plex...")

        # Connect to Plex server without authentication
        plex = PlexServer(PLEX_URL)

        library_data = {}
        for section in plex.library.sections():
            if not ALLOWED_LIBRARIES or section.title in ALLOWED_LIBRARIES:
                logger.info("Processing library: %s", section.title)
                items = []
                for item in section.all():
                    try:
                        thumb_url = None
                        if hasattr(item, "thumbUrl"):
                            # Convert internal URL to external URL
                            thumb_url = (
                                f"{item.thumbUrl}" if item.thumbUrl else None
                            )

                        items.append(
                            {
                                "title": item.title,
                                "year": getattr(item, "year", None),
                                "thumb_url": thumb_url,
                                "added_at": getattr(
                                    item, "addedAt", datetime.now()
                                ).strftime("%Y-%m-%d"),
                                "rating": getattr(item, "rating", None),
                            }
                        )
                    except Exception as e:
                        logger.error("Error processing item %s: %s", item.title, str(e))

                library_data[section.title] = items

        # Save to cache file
        CACHE_FILE.parent.mkdir(parents=True, exist_ok=True)
        with open(CACHE_FILE, "w", encoding="utf-8") as f:
            json.dump(library_data, f)

        return library_data

    except Exception as e:
        logger.error("Error fetching library data: %s", str(e))

        # Try to load from cache file if available
        if CACHE_FILE.exists():
            logger.info("Loading from cache file...")
            with open(CACHE_FILE, "r", encoding="utf-8") as f:
                return json.load(f)
        return {}


async def refresh_cache():
    """Background task to refresh cache periodically"""
    while True:
        try:
            data = await fetch_library_data()
            await cache.set("library_data", data, ttl=REFRESH_INTERVAL)
            logger.info("Cache refreshed successfully")
        except Exception as e:
            logger.error("Error refreshing cache: %s", str(e))
        await asyncio.sleep(REFRESH_INTERVAL)


@app.on_event("startup")
async def startup_event():
    """Initialize cache and start background refresh task"""
    # Initial fetch
    data = await fetch_library_data()
    await cache.set("library_data", data, ttl=REFRESH_INTERVAL)

    # Start background refresh task
    asyncio.create_task(refresh_cache())


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy"}


@app.get("/", response_class=HTMLResponse)
async def home(request: Request):
    """Render home page with library data"""
    try:
        data = await cache.get("library_data")
        if not data:
            data = await fetch_library_data()
            await cache.set("library_data", data, ttl=REFRESH_INTERVAL)

        return templates.TemplateResponse(
            "index.html",
            {
                "request": request,
                "libraries": data,
                "last_updated": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            },
        )
    except Exception as e:
        logger.error("Error rendering home page: %s", str(e))
        return templates.TemplateResponse(
            "error.html", {"request": request, "error": str(e)}
        )


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)
