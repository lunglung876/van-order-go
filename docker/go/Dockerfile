FROM golang:1.12.7

RUN mkdir /go-dep-install
RUN cd /go-dep-install && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/van-order-go
COPY . .

RUN dep ensure
RUN go install van-order-go

EXPOSE 8080
ENTRYPOINT	["/go/bin/van-order-go"]