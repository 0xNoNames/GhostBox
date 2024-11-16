# GhostBox

> Based on [GhostYgg](https://github.com/GuillaumeMCK/GhostYgg)

This application will automatically download any torrent files placed in a designated watched directory, which can be specified via a command-line option or a Docker volume.

The torrent client is configured to avoid sending any information to the torrent tracker, ensuring that the download ratio remains unaffected.

## Installation

### From Binary

You can download the pre-built binaries for your platform from the [releases](https://github.com/0xNoNames/GhostBox/releases) page.

### Using Go

If you are using Go v1.20 or higher, you can install GhostBox using the following command:

```bash
go install -v github.com/0xNoNames/GhostBox@latest
```

### From Source

To get started with GhostBox, follow these steps:

```bash
git clone https://github.com/0xNoNames/GhostBox.git
cd GhostBox
go build GhostBox.go
./GhostBox
```

## Usage

> [!CAUTION]
> You are responsible for the torrents you download with GhostBox.

### Docker CLI

```sh
docker run --rm \
           --name ghostbox \
           -e PUID=501 \
           -e PGID=501 \
           -v ./downloads:/app/downloads \
           -v ./torrents:/app/torrents \
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
      - ./downloads:/app/downloads
      - ./torrents:/app/torrents
    restart: unless-stopped
```

### Using the binary

```sh
./ghostbox -i ./torrents -o ./downloads
```

With the following options:

- `-i`: Specifies the watched directory, where the ".torrent" files will be added.
- `-o`: Specifies the output directory, where the downloaded files will be stored.
- `-help`: Displays the help message.
