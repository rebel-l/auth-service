#!/bin/bash

# stop
docker stop redis
docker rm redis

# start
docker run --name redis -d redis:6-alpine
