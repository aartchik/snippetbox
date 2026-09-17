#!/bin/sh
set -eu

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <image>" >&2
    exit 2
fi

deploy_dir=/home/deploy/apps/snippetbox
mkdir -p "$deploy_dir"
tar -xzf /tmp/snippetbox-deploy.tar.gz -C "$deploy_dir"
rm -f /tmp/snippetbox-deploy.tar.gz
cd "$deploy_dir"
chmod +x deploy/deploy.sh

if [ ! -f .env ]; then
    echo "$deploy_dir/.env is missing; production secrets must be created before deployment" >&2
    exit 1
fi

exec ./deploy/deploy.sh "$1"
