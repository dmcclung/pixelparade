FROM alpine

RUN apk add --update hwloc

RUN mkdir /app

WORKDIR /app

RUN wget https://github.com/filecoin-project/lotus/releases/download/v1.22.1/lotus_v1.22.1_linux_amd64.tar.gz

RUN tar -xvf lotus_1.22.1_linux_amd64.tar.gz

RUN mv lotus_1.22.1_linux_amd64/lotus /usr/local/bin/lotus

ENV FULLNODE_API_INFO=wss://api.chain.love

CMD ["lotus", "daemon", "--lite"]
