#!/bin/bash

source cmd/scripts/base.env
source cmd/scripts/.env
# shellcheck disable=SC2046
#eval $(go run github.com/juxuny/data-utils/cmd/ad-email set-proxy --user=${PROXY_USER} --pass=${PROXY_PASS})
go run github.com/juxuny/data-utils/cmd/split-video-by-subtitle \
  split -i '/Users/juxuny/Downloads/Kung.Fu.Panda.2008.1080p.BluRay.x264.DTS-FGT/KF_2.mp4' \
  --in-srt 'tmp/eng.srt' \
  --out-srt 'tmp/eng.converted.srt' \
  --size=14 \
  --global-size=12 \
  --bg tmp/daily_2.png \
  --cover-color '#0F9883' \
  --cover-duration 2 \
  --cover-font-size 150 \
  --desc-font-color '#393939' \
  --desc-font-size 80 \
  --repeat=3 \
  --expand=1 \
  --max-words=3 \
  "$@"
