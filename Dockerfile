FROM golang:latest
ENV GOPATH=/app
ENV PATH=$GOPATH/bin:$PATH
WORKDIR /app/src/github.com/dnerochain/dnero
COPY . .
RUN make install
RUN cp -r ./integration/testnet_amber ../
EXPOSE 28888
CMD dnero start --config=../testnet_amber/node --password="qwertyuiop"

