FROM golang:1.19

RUN useradd gofer
USER gofer
WORKDIR /home/gofer

ENV CGO_ENABLED=0

COPY ./go.mod /home/gofer/src/goferos/
COPY ./go.sum /home/gofer/src/goferos/
COPY ./cmd/ /home/gofer/src/goferos/cmd/
RUN cd /home/gofer/src/goferos \
    && go install ./cmd/wget \
    && go install ./cmd/gosh \
    && go install ./cmd/git \
    ;

RUN wget -O /tmp/getmicro.sh https://getmic.ro/
RUN cd /tmp && gosh getmicro.sh && mv ./micro /go/bin

CMD ["gosh"]
