#syntax=docker.io/docker/dockerfile:1.2

FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS build
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
RUN --mount=target=/mnt ["/mnt/build.sh"]

FROM scratch
COPY --from=build /icmp-tunnel /
ENTRYPOINT ["/icmp-tunnel"]
