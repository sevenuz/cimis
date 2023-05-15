# syntax=docker/dockerfile:1

## Build node
FROM node:18.13 AS build_frontend

WORKDIR /ui

COPY ./ui ./
RUN npm install

RUN npm run build

## Build go
FROM golang:1.18-buster AS build_backend

WORKDIR /

COPY --from=build_frontend /ui/build /ui/build
COPY ./migrations ./migrations
COPY ./ui/embed.go ./ui/embed.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./main.go ./main.go

RUN go mod download

RUN go build

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build_backend /cimis /cimis

EXPOSE 8090

ENTRYPOINT ["/cimis", "serve"]

