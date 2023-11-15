FROM alpine:3.16.7

ARG LOTUS_VERSION=v1.23.3
ARG FULLNODE_API_INFO=wss://api.chain.love

RUN apk add --no-cache hwloc \
    && wget https://github.com/filecoin-project/lotus/releases/download/${LOTUS_VERSION}/lotus_${LOTUS_VERSION}_linux_amd64.tar.gz -O /tmp/lotus.tar.gz \
    && tar -xzf /tmp/lotus.tar.gz -C /tmp \
    && mv /tmp/lotus_${LOTUS_VERSION}_linux_amd64/lotus /usr/local/bin/lotus \
    && rm -rf /tmp/*

WORKDIR /usr/local/bin

ENV FULLNODE_API_INFO=${FULLNODE_API_INFO}

ENTRYPOINT ["lotus"]

CMD ["daemon", "--lite"]
