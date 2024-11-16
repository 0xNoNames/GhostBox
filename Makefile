.DEFAULT_GOAL := default

IMAGE ?= xnonames/ghostbox:latest

export DOCKER_CLI_EXPERIMENTAL=enabled

.PHONY: build # Build the container image
build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--output "type=docker,push=false" \
		--tag $(IMAGE) \
		.

.PHONY: publish # Push the image to the remote registry
publish:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--platform linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm/v8,linux/arm64,linux/ppc64le,linux/s390x \
		--output "type=image,push=true" \
		--tag $(IMAGE) \
		.


.PHONY: generate
generate:
		CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o dist/ghostbox_linux_i386 GhostBox.go              # Linux 32bit
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/ghostbox_linux_x86_64 GhostBox.go          # Linux 64bit
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o dist/ghostbox_linux_arm GhostBox.go       # Linux armv5/armel/arm (it also works on armv6)
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o dist/ghostbox_linux_armhf GhostBox.go     # Linux armv7/armhf
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/ghostbox_linux_aarch64 GhostBox.go         # Linux armv8/aarch64
		CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o dist/ghostbox_freebsd_x86_64 GhostBox.go      # FreeBSD 64bit
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/ghostbox_darwin_x86_64 GhostBox.go        # Darwin 64bit
		CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/ghostbox_darwin_aarch64 GhostBox.go       # Darwin armv8/aarch64
		CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o dist/ghostbox_windows_i386.exe GhostBox.go      # Windows 32bit
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/ghostbox_windows_x86_64.exe GhostBox.go  # Windows 64bit
