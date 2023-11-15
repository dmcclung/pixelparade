FROM ubuntu:22.04

ARG LOTUS_VERSION=v1.23.3
ARG FULLNODE_API_INFO=wss://api.chain.love

RUN apt-get update \
    && apt-get install -y hwloc wget \
    && wget https://github.com/filecoin-project/lotus/releases/download/${LOTUS_VERSION}/lotus_${LOTUS_VERSION}_linux_amd64.tar.gz -O /tmp/lotus.tar.gz \
    && tar -xzf /tmp/lotus.tar.gz -C /tmp \
    && mv /tmp/lotus_${LOTUS_VERSION}_linux_amd64/lotus /usr/local/bin/lotus \
    && chown root:root /usr/local/bin/lotus \
    && rm -rf /tmp/*

WORKDIR /usr/local/bin

ENV FULLNODE_API_INFO=${FULLNODE_API_INFO}

ENTRYPOINT ["lotus"]

CMD ["daemon", "--lite"]
