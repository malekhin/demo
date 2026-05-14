#!/usr/bin/env bash

migrate_wr() {
    /migrate_pg.sh down
    /migrate_pg.sh up
}
migrate_wr
while [ $? -ne 0 ]; do
    sleep 1
    migrate_wr
done

go build -buildvcs=false -o ./tmp/app ./src/cmd/app

air
