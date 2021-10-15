#!/bin/bash

source cmd/scripts/base.env
source cmd/scripts/.env
# shellcheck disable=SC2046
#eval $(go run github.com/juxuny/data-utils/cmd/ad-email set-proxy --user=${PROXY_USER} --pass=${PROXY_PASS})
go run github.com/juxuny/data-utils/cmd/split-video-by-subtitle "$@"
