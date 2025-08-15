#!/bin/sh

set -e

# to load env before running db migration
source /app/app.env

echo "start the app"
exec "$@"