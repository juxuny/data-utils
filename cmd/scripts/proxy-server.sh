#!/bin/bash
source cmd/scripts/base.env
go run "${PACKAGE_NAME}/proxy-server" -c tmp/proxy.yaml