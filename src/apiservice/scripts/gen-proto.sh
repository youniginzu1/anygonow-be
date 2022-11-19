#!/bin/bash
#
# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
GOPATH=$HOME/go
PATH=$PATH:$GOPATH/bin
protodir=../../proto

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc --go_out ./src/pb --go_opt paths=source_relative \
   --go-grpc_out ./src/pb --go-grpc_opt paths=source_relative \
   --validate_out="lang=go:./src" \
   --grpc-gateway_out ./src/pb --grpc-gateway_opt paths=source_relative --grpc-gateway_opt allow_delete_body=true \
   --openapiv2_out ./src/services/swagger --openapiv2_opt logtostderr=true --openapiv2_opt allow_delete_body=true \
   -I $protodir -I=$GOPATH/src $protodir/apiservice.proto

protoc --go_out ./src/pb/mailpb --go_opt paths=source_relative \
   --go-grpc_out ./src/pb/mailpb --go-grpc_opt paths=source_relative \
   -I $protodir $protodir/mailservice.proto

protoc --go_out ./src/pb/chatpb --go_opt paths=source_relative \
   --go-grpc_out ./src/pb/chatpb --go-grpc_opt paths=source_relative \
   -I $protodir $protodir/chatservice.proto

protoc --go_out ./src/pb/authpb --go_opt paths=source_relative \
   --go-grpc_out ./src/pb/authpb --go-grpc_opt paths=source_relative \
   -I $protodir $protodir/authservice.proto

protoc --go_out ./src/internal/var/c --go_opt paths=source_relative \
   -I $protodir $protodir/const.proto