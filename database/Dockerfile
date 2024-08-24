FROM postgres:16-bookworm

RUN apt-get update \
    && apt-get install -y \
    git curl wget unzip build-essential postgresql-server-dev-16 \
    libssl-dev libcrypto++-dev zlib1g-dev python3 python3-pip pkg-config \
    libreadline-dev liblz4-dev \
    && python3 -m pip install --upgrade pip
