# Copyright (c) 2021 jblim Mobigen, Inc. All Rights Reserved.

# Build the docker
ARG DOCKER_BUILD_IMAGE=golang:1.26.1-alpine3.23
ARG DOCKER_BASE_IMAGE=alpine:3.12

FROM ${DOCKER_BUILD_IMAGE} AS build
WORKDIR /app/
COPY . /app/
RUN apk update && apk add --no-cache make build-base
RUN make build

# Final Image
FROM ${DOCKER_BASE_IMAGE}

LABEL name="web-server" \
  maintainer="irisdev@mobigen.com" \
  vendor="Mobigen" \
  architecture="x86_64" \
  url="https://mobigen.com" \
  description="Web Server Sample. Created Platform Team Mobigen"

#ENV SERVER=/app/bin/server \
ENV USER_UID=1000 \
    USER_NAME=mobigen \
    APP_HOME=/app \
    PROFILE=prod \
    LOG_LEVEL=info

RUN  apk update && apk add --no-cache libc6-compat ca-certificates
COPY --from=build /app/scripts/user_setup /usr/local/bin/user_setup
COPY --from=build /app/build/bin /app/bin
COPY --from=build /app/configs /app/configs

RUN chmod +x /usr/local/bin/user_setup && /usr/local/bin/user_setup

WORKDIR /app

USER ${USER_UID}

EXPOSE 8080

ENTRYPOINT ["/app/bin/test"]
