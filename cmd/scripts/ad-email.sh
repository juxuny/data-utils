#!/bin/bash

source cmd/scripts/base.env
source cmd/scripts/.env
# shellcheck disable=SC2046
#eval $(go run github.com/juxuny/data-utils/cmd/ad-email set-proxy --user=${PROXY_USER} --pass=${PROXY_PASS})
go run github.com/juxuny/data-utils/cmd/ad-email send \
  --send-file=tmp/sender.list \
  --verbose \
  --db-host=127.0.0.1 \
  --db-port=3307 \
  --db-user=root \
  --db-pwd=123456 \
  --db-name=crawl \
  --host=smtp.163.com:465 \
  --ssl --batch-size=20 --delay=60 --count=1