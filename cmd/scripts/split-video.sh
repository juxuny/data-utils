#!/bin/bash

source cmd/scripts/base.env
source cmd/scripts/.env
# shellcheck disable=SC2046
#eval $(go run github.com/juxuny/data-utils/cmd/ad-email set-proxy --user=${PROXY_USER} --pass=${PROXY_PASS})
go run github.com/juxuny/data-utils/cmd/split-video-by-subtitle \
  split -i '/Users/juxuny/Downloads/Kung.Fu.Panda.2008.1080p.BluRay.x264.DTS-FGT/out.mp4' \
  --bg tmp/daily_2.png \
  --cover-color '#ea5a4f' \
  --cover-duration 2 \
  --cover-font-size 120 \
  --desc-font-color '#393939' \
  --desc-font-size 96 \
  "$@"
