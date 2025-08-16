#!/bin/sh

set -e

echo "run db migration"
source /app/app.env

echo "start the app"
exec "$@"