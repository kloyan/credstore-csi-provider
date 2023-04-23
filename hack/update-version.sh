#!/usr/bin/env bash

set -euo pipefail

VERSION="$1"

sed -i \
    -e "s|VERSION?=\(.*\)|VERSION?=$VERSION|g" Makefile \
    -e "s|credstore-csi-provider:\(.*\)|credstore-csi-provider:$VERSION|g" deploy/daemonset.yaml
