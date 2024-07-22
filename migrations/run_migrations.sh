#!/bin/sh
set -e

echo "Waiting for database to be ready..."
/root/wait-for-it.sh $DB_HOST:5432 -t 60

echo "Running database migrations..."
/usr/local/bin/migrate -path /root/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:5432/$DB_NAME?sslmode=disable" up

echo "Migrations completed."

exec "$@"