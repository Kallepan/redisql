FROM golang:1.21.2-bookworm

LABEL name="GoLang Dev Container" \
    description="GoLang Dev Container" \
    version="1.0.0"

# Install tools
RUN apt-get update && apt-get install -y \
    curl \
    git \
    gnupg2 \
    lsb-release \
    sudo \
    unzip \
    wget \
    zsh \
    && rm -rf /var/lib/apt/lists/*