#!/bin/sh
set -eu

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <image>" >&2
    exit 2
fi

export APP_IMAGE="$1"
compose="docker compose --project-name snippetbox --env-file .env -f deploy/compose.prod.yml"

$compose pull db redis migrate caddy
if ! $compose pull app; then
    docker image inspect "$APP_IMAGE" >/dev/null
fi
$compose up -d db redis
$compose run --rm migrate
$compose up -d --remove-orphans app caddy

i=0
until curl --fail --silent --show-error http://127.0.0.1:4000/ping >/dev/null; do
    i=$((i + 1))
    if [ "$i" -ge 30 ]; then
        $compose logs --tail=100 app
        exit 1
    fi
    sleep 2
done

$compose ps
