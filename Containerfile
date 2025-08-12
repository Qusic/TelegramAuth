FROM --platform=$BUILDPLATFORM docker.io/library/golang:alpine AS build
ARG TARGETOS TARGETARCH
ARG GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN go build

FROM scratch
COPY --from=build /app/TelegramAuth /
ENTRYPOINT [ "/TelegramAuth" ]
