FROM golang:1.16-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /app/TelegramAuth /
WORKDIR /data
ENTRYPOINT [ "/TelegramAuth" ]
