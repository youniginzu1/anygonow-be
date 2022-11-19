#!/bin/sh

go install github.com/google/wire/cmd/wire@latest
wire ./...