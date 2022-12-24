FROM ubuntu:22.04

ENV NODE_MAJOR_VERSION=16
ENV GOLANG_VERSION=1.19.4

# utilities
RUN apt update
RUN apt install -y curl jq git

# golang
RUN curl --silent --location https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz --output go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"

# nodejs
RUN curl --silent --location https://deb.nodesource.com/setup_${NODE_MAJOR_VERSION}.x | bash -
RUN apt install -y nodejs
