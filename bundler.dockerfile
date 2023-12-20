# Use the official Golang image to create a build artifact.
FROM golang:1.20-alpine as builder

# Set up apk dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev curl ca-certificates build-base

# Set working directory for the build
WORKDIR /opt/app

# Add source files
COPY . .

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache $PACKAGES

# For Private REPO
ARG GH_TOKEN=""
RUN go env -w GOPRIVATE="github.com/node-real/*"
RUN git config --global url."https://${GH_TOKEN}@github.com".insteadOf "https://github.com"

# Build the bundler binary using makefile
RUN make build-bundler

# Use the Alpine-based runtime image to create a lean production container.
FROM alpine:3.17

ARG USER=app
ARG USER_UID=1000
ARG USER_GID=1000

ENV PACKAGES ca-certificates libstdc++ curl
ENV WORKDIR=/app

RUN apk add --no-cache $PACKAGES \
  && rm -rf /var/cache/apk/* \
  && addgroup -g ${USER_GID} ${USER} \
  && adduser -u ${USER_UID} -G ${USER} --shell /sbin/nologin --no-create-home -D ${USER} \
  && addgroup ${USER} tty \
  && sed -i -e "s/bin\/sh/bin\/bash/" /etc/passwd

WORKDIR ${WORKDIR}
RUN chown -R ${USER_UID}:${USER_GID} ${WORKDIR}
USER ${USER_UID}:${USER_GID}

ENV CONFIG_FILE_PATH /opt/app/config/config.json

# Copy the binary to the production image from the builder stage.
COPY --from=builder /opt/app/build/bundler /app/bundler

# Run the bundler service on container startup.
CMD ["/app/bundler", "--config-path", "$CONFIG_FILE_PATH"]