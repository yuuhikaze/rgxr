#!/usr/bin/env bash

set -uo pipefail
IFS=$'\n'

docker compose down
docker volume rm -f rgxr_pgdata
