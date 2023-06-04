# syntax=docker/dockerfile:1

################################
# Build
################################
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN make build

################################
# Deploy
################################
FROM --platform=linux/amd64 debian:bullseye-slim

# Install Dependencies
RUN apt-get update  \
    && apt-get install -y ca-certificates wget xz-utils curl unzip python git \
    && update-ca-certificates

# Copy built binaries
RUN mkdir -p /opt/posts-api
COPY --from=build /app/bin/posts-api /opt/posts-api/
COPY --from=build /app/deploy/local/local.env /opt/posts-api/service.env
COPY --from=build /app/scripts /opt/posts-api/scripts
COPY --from=build /app/migrations /app/migrations
WORKDIR /opt/posts-api

# Run app
EXPOSE 8080
ENTRYPOINT ["./scripts/entrypoint.sh"]