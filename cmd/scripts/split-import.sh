#!/bin/bash

source cmd/scripts/base.env
source cmd/scripts/.env
# shellcheck disable=SC2046
go run github.com/juxuny/data-utils/cmd/split-video-by-subtitle \
  import \
  --srt-dir=./tmp/subtitle \
  "$@"
