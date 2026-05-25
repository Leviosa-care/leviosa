FROM node:22-bookworm

# System dependencies
RUN apt-get update && apt-get install -y \
  git \
  curl \
  jq \
  && rm -rf /var/lib/apt/lists/*

# Go 1.24.2 — matches backend go.mod requirement
RUN curl -fsSL https://go.dev/dl/go1.24.2.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:$PATH"

# UID/GID alignment: sandcastle docker build-image passes host UID/GID automatically
ARG AGENT_UID=1000
ARG AGENT_GID=1000

RUN groupmod -g $AGENT_GID node && \
    usermod -u $AGENT_UID -g $AGENT_GID -d /home/agent -m -l agent node

USER ${AGENT_UID}:${AGENT_GID}

ENV GOPATH="/home/agent/go"
ENV PATH="/home/agent/go/bin:/home/agent/.local/bin:$PATH"

# Claude Code CLI (used by both GLM and Anthropic agents — same binary, different env)
RUN curl -fsSL https://claude.ai/install.sh | bash

WORKDIR /home/agent
ENTRYPOINT ["sleep", "infinity"]
