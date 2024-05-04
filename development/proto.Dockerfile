
FROM golang:1.22.1 AS protoc_gen_go

RUN apt update && apt install -y --no-install-recommends curl make git unzip apt-utils
ENV GO111MODULE=on
ENV PROTOC_VERSION=3.14.0
ENV GRPC_WEB_VERSION=1.2.1
ENV BUFBUILD_VERSION=0.24.0

RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip
RUN unzip protoc-$PROTOC_VERSION-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/envoyproxy/protoc-gen-validate@latest
RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

WORKDIR /app
COPY .. .

RUN go install ./protoc/protoc-gen-event/.