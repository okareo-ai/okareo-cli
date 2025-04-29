FROM --platform=$BUILDPLATFORM python:3.9-slim

WORKDIR /app

ARG TARGETARCH

RUN apt-get update && apt-get install -y curl && \
    if [ "$TARGETARCH" = "arm64" ]; then \
        curl -O -L https://github.com/okareo-ai/okareo-cli/releases/latest/download/okareo_linux_arm64.tar.gz && \
        tar -xvf okareo_linux_arm64.tar.gz && \
        rm okareo_linux_arm64.tar.gz; \
    else \
        curl -O -L https://github.com/okareo-ai/okareo-cli/releases/latest/download/okareo_linux_amd64.tar.gz && \
        tar -xvf okareo_linux_amd64.tar.gz && \
        rm okareo_linux_amd64.tar.gz; \
    fi && \
    chmod +x okareo && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

ENV PATH="/app:${PATH}"
EXPOSE 4000

CMD ["okareo", "proxy"]
