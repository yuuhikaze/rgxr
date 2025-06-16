#!/usr/bin/env sh
set -e

psql -v ON_ERROR_STOP=1 \
    --username "$POSTGRES_USER" \
    --dbname "$POSTGRES_DB" \
    --set=auth_pass="$PG_AUTHENTICATOR_PASSWORD" <<- EOSQL
    create role web_anon nologin;
    grant usage on schema api to web_anon;
    grant select on api.finite_automatas to web_anon;

    create role authenticator noinherit login password :'auth_pass';
    grant web_anon to authenticator;
EOSQL
