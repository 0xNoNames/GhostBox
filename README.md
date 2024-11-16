# GhostBox

> Based on [GhostYgg](https://github.com/GuillaumeMCK/GhostYgg)

This application will automatically download any torrent files placed in a designated watched directory, which can be specified via a command-line option or a Docker volume.

The torrent client is configured to avoid sending any information to the torrent tracker, ensuring that the download ratio remains unaffected.

## Usage

### Docker CLI

```sh
docker run --rm \
           --name ghostbox \
           -e PUID=501 \
           -e PGID=501 \
           -v ./downloads:/downloads \
           -v ./torrents:/torrents \
          ghostbox:latest
```

### Docker Compose

```yaml
---
services:
  ghostbox:
    image: ghostbox:latest
    container_name: ghostbox
    environment:
      - PUID=501
      - PGID=501
    volumes:
      - ./downloads:/downloads
      - ./torrents:/torrents
    restart: unless-stopped
```

### Using directly the binary

```sh
./ghostbox -i ./torrents -o ./downloads
```

