#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migrations -database "$POSTGRES_URL" -verbose up

echo "start the app"
exec "$@"
