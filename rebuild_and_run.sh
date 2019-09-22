#!/bin/sh

cd public/ \
    && ./build.sh \
    && cd ../ \
    && go run server.go