#!/bin/bash

sudo docker rm postgres --force
sudo docker run -d -p5432:5432 -v `pwd`/postgres-data:/var/lib/postgresql/data -e POSTGRES_USER="tapit" -e POSTGRES_PASSWORD="secret-tapit-password" -e POSTGRES_DB="tapit" --name postgres postgres
