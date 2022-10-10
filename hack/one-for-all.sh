#!/bin/bash

make -C src/apiservice dev-recreate
make -C src/chatservice dev-recreate
make -C src/authservice dev-recreate
make -C src/mailservice dev-recreate
make -C src/loadbalancer dev-recreate