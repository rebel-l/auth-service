#!/bin/bash

BRANCH=`git rev-parse --abbrev-ref HEAD | sed -r 's/\/+/-/g'`

# stop
sudo docker stop auth-service
sudo docker rm auth-service

# build
sudo docker build -t rebel1l/auth-service:$BRANCH .

# start
sudo docker run --name auth-service -d -it -p 3000:3000 rebel1l/auth-service:$BRANCH
sudo docker logs auth-service
sudo docker ps
