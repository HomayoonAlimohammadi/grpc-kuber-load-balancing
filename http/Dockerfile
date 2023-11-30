FROM golang:1.21 AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /api ./post_api

FROM alpine:latest AS release

WORKDIR /

COPY --from=build /api /api

ENTRYPOINT ["/api"]
