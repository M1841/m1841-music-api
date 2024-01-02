# syntax=docker/dockerfile:1

FROM golang:bullseye AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /m1841-music-api

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /m1841-music-api /m1841-music-api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/m1841-music-api" ]