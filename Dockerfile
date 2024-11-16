FROM golang:1.23.3 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o ghostbox GhostBox.go
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ghostbox GhostBox.go


# Use distroless as minimal base image to package the manager binary
FROM alpine:latest
# Create the /torrents and /downloads directories
WORKDIR /app
RUN mkdir torrents downloads
COPY --from=builder /workspace/ghostbox .
CMD [ "/app/ghostbox", "-i", "/app/torrents", "-o", "/app/downloads" ]
