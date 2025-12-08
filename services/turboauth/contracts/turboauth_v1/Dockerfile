# Multi-stage build for C++ smart contract
FROM ubuntu:22.04 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /contracts

# Copy contract source
COPY . .

# Build contract
RUN mkdir -p build && cd build && \
    cmake .. && \
    make

# Final stage - minimal image with compiled contract
FROM scratch

# Copy compiled contract
COPY --from=builder /contracts/build/libmicroauth.a /contract/

# Metadata
LABEL description="Qubic MicroAuth Smart Contract v1"
LABEL version="1.0.0"
