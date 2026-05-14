#!/usr/bin/env bash

migrate -database "postgres://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/$PG_DB?sslmode=disable" -path /app/migrations $@
