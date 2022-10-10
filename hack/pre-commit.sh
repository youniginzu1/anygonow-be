#!/bin/bash

make -C src/apiservice proto
make -C src/authservice proto
make -C src/chatservice proto
make -C src/mailservice proto

make -C src/apiservice test