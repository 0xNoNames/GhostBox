# Specifies a parent image
FROM golang:latest

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Installs Go dependencies
RUN go mod download

# Builds your app with optional configuration
RUN go build -o /ghostbox

# Create the /torrents and /downloads directories
RUN mkdir /torrents /downloads

# Specifies the executable command that runs when the container starts
CMD [ "/ghostbox", "-i", "/torrents", "-o", "/downloads" ]
