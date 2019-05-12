#!/usr/bin/env bash

set -euo pipefail

godoc2md github.com/andy2046/rund \
    > $GOPATH/src/github.com/andy2046/rund/docs.md
