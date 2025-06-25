#!/usr/bin/env bash

set -uo pipefail
IFS=$'\n'

curl -f -s 'http://localhost/pgapi/finite_automatas' > /dev/null || {
    echo 'Postgrest failed!'
    return 1
}
curl -f -s 'http://localhost/api/complement?uuid=2dd3db92-6834-41fe-9190-6cf6354c04c5' > /dev/null || {
    echo 'Backend failed!'
    return 1
}
echo "Postgrest and backend are healthy!"
