FROM golang:1.14-alpine AS build
WORKDIR /app
ENV CGO_ENABLED=0
COPY . .
RUN go build

FROM scratch
COPY --from=build /app/TelegramAuth /
ENTRYPOINT [ "/TelegramAuth" ]
