# Plex Kiosk

![alt text](image.png)

> [!NOTE]  
> These are just thumbnails and do not demonstrate nor imply ownership or possession of any content

Grab your libraries and display items in a kiosk view. Useful as a scrollable view or sharing your library with friends. Click on an item to show a description, rating, and other associated metadata.

This is WIP. I dont know what I want to do with this yet. I made this in about an hour. I might add movie poster functionality (sync with now playing and/or display a watchlist, maybe support user requests so your friends can vote on something to watch async and it will show up)

## Usage

### Plex Authentication
For now, you will have to set Plex to allow local access without authentication (add your server IP to the list of allowed networks).

I might implement Plex OAuth in the future but its not ideal because you would need to constantly reuathenticate. Plex really needs to implement local only API keys. Plex tokens are not ideal because they also expire.

### Unraid/Orchestrator

Create a new container with the following settings:

- Repository: `ghcr.io/iloveicedgreentea/plex-kiosk:latest`
- Environment Variables:
  - `PLEX_URL`: `scheme://yourplexserver:port`
  - `ALLOWED_LIBRARIES`: `Movies,TV Shows` # comma separated list of libraries names
  - `REFRESH_INTERVAL`: `21600` # 6 hours or whatever you want
  - `TZ`: `UTC` # if you want
- Ports:
    - `PORT:80` # map whatever you want to container port 80

I personally use cloudflare tunnels to expose this to the internet. You can set oauth or password login that way.

### Local Testing

#### Docker Compose
`just run-docker`

#### Backend

`export PLEX_URL="scheme://yourplexserver:port"`
`export ALLOWED_LIBRARIES="Movies,TV Shows"`

`just run`

#### Frontend
`cd frontend`

`bun run dev`