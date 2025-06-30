#!/usr/bin/env bash

set -uo pipefail
IFS=$'\n'

docker exec -u postgres rgxr-postgres-1 \
  pg_dump -d postgres -t api.finite_automatas -f /tmp/finite_automatas.dump
docker cp rgxr-postgres-1:/tmp/finite_automatas.dump ../finite_automatas.dump
