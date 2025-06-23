#!/usr/bin/env bash

set -uo pipefail
IFS=$'\n'

docker cp "$1":/data ../data
