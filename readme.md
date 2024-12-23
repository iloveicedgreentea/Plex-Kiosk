# Plex Kiosk

This is WIP

Grabs your libraries and displays it in a new kiosk view. Useful as a scrollable view or sharing your library with friends.


## Usage

For now, you will have to set Plex to allow local access without authentication (add your server IP to the list of allowed networks).

I might implement Plex OAuth in the future.

### Local
`export PLEXURL="scheme://yourplexserver:port"`

refer to docker compose for other env vars

`just run`

### Unraid/Orchestrator
Create a new container with the following settings:

- Repository: `ghcr.io/tbd`
- Environment Variables:
  - `PLEXURL`: `scheme://yourplexserver:port`
  - `ALLOWED_LIBRARIES`: `Movies,TV Shows` 
  - `REFRESH_INTERVAL`: `21600` # 6 hours
  - `TZ`: `UTC` # if you want
- Ports:
    - `PORT:80` # map whatever you want

I personally use cloudflare tunnels to expose this to the internet.