# Build binary
FROM --platform=$BUILDPLATFORM golang:1.23.3-alpine AS build-env
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
# Copy the go sources
COPY . .
# Build
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -ldflags="-s -w" -o ghostbox GhostBox.go

RUN mkdir torrents downloads


# Create image
# Use distroless as minimal base image to package the manager binary
FROM scratch
WORKDIR /app
COPY --from=build-env /workspace/ghostbox .
COPY --from=build-env /workspace/torrents .
COPY --from=build-env /workspace/downloads .
CMD [ "/app/ghostbox", "-i", "/app/torrents", "-o", "/app/downloads" ]
