#!/bin/sh

#* variables
PROTO_PATH=./proto
PROTO_OUT=./pb
IDL_PATH=./
DOC_OUT=./docs

#* clean
# rm -rf ${PROTO_OUT}
protoc \
  ${PROTO_PATH}/**/*.proto \
  -I=/usr/local/include \
  --proto_path=${PROTO_PATH} \
  --go_out=:${IDL_PATH} \
  --validate_out=lang=go:${IDL_PATH} \
  --go-grpc_out=:${IDL_PATH} \
  --grpc-gateway_out=:${IDL_PATH} \
  --event_out=:${IDL_PATH} \
  --enum_out=:${IDL_PATH} \
  --http_out=:${IDL_PATH} \
  --openapiv2_out=:${DOC_OUT}/swagger
