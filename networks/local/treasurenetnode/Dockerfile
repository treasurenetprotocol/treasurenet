# FROM golang:stretch as build-env
FROM golang:buster as build-env
# FROM golang:1.22.4-alpine3.20 AS build-env
# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

# Set working directory for the build
# WORKDIR /go/src/github.com/treasurenet
WORKDIR /root/go/src/treasurenet
# Add source files
COPY . .

# build Ethermint
RUN make build-linux

# Final image
FROM golang:1.18 as final
# FROM golang:1.22.4 as final

WORKDIR /

RUN apt-get update

# Copy over binaries from the build-env
# COPY --from=build-env /go/src/github.com/treasurenet/build/treasurennetd /
# COPY --from=build-env /go/src/github.com/treasurenet/scripts/start-docker.sh /
COPY --from=build-env /root/go/src/treasurenet/build/treasurenetd /
COPY --from=build-env /root/go/src/treasurenet/scripts/start-docker.sh /

EXPOSE 26656 26657 1317 8545 8546

# Run ethermintd by default, omit entrypoint to ease using container with treasurenetd
ENTRYPOINT ["/bin/bash", "-c"]