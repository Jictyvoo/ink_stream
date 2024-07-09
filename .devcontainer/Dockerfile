FROM ghcr.io/ciromattia/kcc:6.1.0

# Run updates and clean up
RUN apt update && apt install -y wget && apt upgrade -y  \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*

# Install Go
RUN wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz \
    && rm go1.22.5.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
USER root
RUN pip3 install --upgrade pip \
    && pip3 install natsort
