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

protoc --go_out ./src/pb/mailpb --go_opt paths=source_relative \
   --go-grpc_out ./src/pb/mailpb --go-grpc_opt paths=source_relative \
   -I $protodir  $protodir/mailservice.proto

protoc --go_out ./src/pb/chatpb --go_opt paths=source_relative \
   --go-grpc_out ./src/pb/chatpb --go-grpc_opt paths=source_relative \
   -I $protodir  $protodir/chatservice.proto

protoc --go_out ./src/internal/var/c --go_opt paths=source_relative \
   -I $protodir $protodir/const.proto