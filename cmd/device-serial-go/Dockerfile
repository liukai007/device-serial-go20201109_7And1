#
# Copyright (C) 2018 IOTech Ltd
#
# SPDX-License-Identifier: Apache-2.0
FROM golang:1.12-alpine AS builder

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --update make git openssh build-base

# set the working directory
WORKDIR $GOPATH/src/github.com/edgexfoundry/device-serial-go

COPY . .

RUN make build
#RUN  make cmd/device-modbus

FROM scratch

ENV APP_PORT=49990
EXPOSE $APP_PORT

COPY --from=builder /go/src/github.com/edgexfoundry/device-serial-go/cmd/device-serial-go  /

LABEL license='SPDX-License-Identifier: Apache-2.0' \
      copyright='Copyright (c) 2019: IoTech Ltd'

ENTRYPOINT ["/device-serial-go","--profile=docker","--confdir=/res","--registry=consul://edgex-core-consul:8500"]