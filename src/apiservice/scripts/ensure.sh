#!/bin/sh
dlv version && echo "==========================>> Dlv ready <<==========================" || go install github.com/go-delve/delve/cmd/dlv@latest 

grpc-client-cli -v && echo "====================>> grpc-client-cli ready <<====================" || go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@v1.10.0
