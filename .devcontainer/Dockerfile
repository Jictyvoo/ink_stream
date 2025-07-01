FROM ghcr.io/ciromattia/kcc:6.2.1

# Run updates and clean up
RUN apt update && apt install -y wget make && apt upgrade -y  \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /home/Downloads/
ENV GOVERSION=1.24.3
ENV COMPRESSED_GO_FILE=go${GOVERSION}.linux-amd64.tar.gz

# Install Go
RUN wget https://go.dev/dl/${COMPRESSED_GO_FILE} \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf ${COMPRESSED_GO_FILE} \
    && rm ${COMPRESSED_GO_FILE}

ENV PATH=$PATH:/usr/local/go/bin
USER root
RUN pip3 install --upgrade pip \
    && pip3 install natsort

WORKDIR /home/kindlego
