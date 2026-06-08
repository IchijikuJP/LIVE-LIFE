FROM golang:1.24-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
COPY cmd ./cmd
RUN go build -o /out/livelife-api ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=build /out/livelife-api /app/livelife-api
COPY static /app/static

ENV PORT=8080
EXPOSE 8080

CMD ["/app/livelife-api"]
