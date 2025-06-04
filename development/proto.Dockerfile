FROM golang:1.24.3 AS install-stage

ENV GOBIN=/usr/local/bin

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/envoyproxy/protoc-gen-validate@v0.9.1
RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

RUN go mod download github.com/googleapis/googleapis@v0.0.0-20221209211743-f7f499371afa

ENV MOD=$GOPATH/pkg/mod
RUN mkdir -p /collect/validate && cp -r $MOD/github.com/envoyproxy/protoc-gen-validate@v0.9.1/validate /collect/validate

RUN cp -r $MOD/github.com/googleapis/googleapis@v0.0.0-20221209211743-f7f499371afa/google /collect/.

WORKDIR /app
COPY . .

RUN go install ./protoc/protoc-gen-event/.
RUN go install ./protoc/protoc-gen-enum/.
RUN go install ./protoc/protoc-gen-api-info/.
RUN go install ./protoc/protoc-gen-http-client/.

FROM alpine:latest AS generate-stage
RUN apk add curl unzip libc6-compat

ENV PROTOC_VERSION=3.14.0
ENV GRPC_WEB_VERSION=1.2.1
ENV BUFBUILD_VERSION=0.24.0

RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip
RUN unzip protoc-$PROTOC_VERSION-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/ /usr/local/include/

# copy binary files
COPY --from=install-stage /usr/local/bin/ /usr/local/bin/
COPY --from=install-stage /bin/protoc-gen-* /usr/local/bin/

# copy default proto files
COPY --from=install-stage /usr/local/include/* /usr/local/include/

COPY --from=install-stage /collect/google/ /usr/local/include/google/
COPY --from=install-stage /collect/validate /usr/local/include/

COPY --from=install-stage /app/proto /proto