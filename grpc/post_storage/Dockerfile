FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /storage

FROM alpine:latest AS release

WORKDIR /

COPY --from=build /storage /storage

ENTRYPOINT ["/storage"]
