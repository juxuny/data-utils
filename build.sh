#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o ./build/split-video-by-subtitle ./cmd/split-video-by-subtitle