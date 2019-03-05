#!/bin/bash
docker run --rm -d --name nats -p 4222:4222 -p 8222:8222 nats-streaming
docker run --rm -d --name mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true -p 3306:3306 mysql
