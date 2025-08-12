#!/usr/bin/env fish
go test -v ./... && staticcheck || exit 1
set manifest ghcr.io/qusic/telegram-auth:latest
podman manifest rm $manifest --ignore
podman manifest create $manifest \
  --annotation org.opencontainers.image.description=TelegramAuth \
  --annotation org.opencontainers.image.source=https://github.com/Qusic/TelegramAuth \
  --annotation org.opencontainers.image.licenses=MIT
podman build . --manifest $manifest \
  --platform linux/amd64 \
  --platform linux/arm64
podman manifest inspect $manifest
podman manifest push $manifest --all --rm
podman image prune --force
