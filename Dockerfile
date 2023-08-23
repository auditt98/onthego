FROM ghcr.io/zitadel/zitadel:latest

RUN apt-get update && apt-get install -y ca-certificates