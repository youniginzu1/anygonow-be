#!/bin/sh
if [[ -z $GCFLAGS ]]
then
    make dev
else
    make debug
fi