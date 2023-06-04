FROM --platform=linux/amd64 debian:bullseye-slim

# Install Dependencies
RUN apt-get update  \
    && apt-get install -y ca-certificates wget xz-utils curl unzip python git \
    && update-ca-certificates

# Copy built binaries
RUN mkdir -p /opt/posts-api
COPY bin/posts-api /opt/posts-api/
COPY scripts /opt/posts-api/scripts
COPY migrations /root/project/posts-api/migrations
WORKDIR /opt/posts-api

# Run app
EXPOSE 8080
ENTRYPOINT ["./scripts/entrypoint.sh"]