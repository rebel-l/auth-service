#!/bin/bash
STORAGE=`pwd`/storage
BRANCH=`git rev-parse --abbrev-ref HEAD | sed -r 's/\/+/-/g'`

echo
echo "Restart Auth Service ..."
echo "Branch: $BRANCH"
echo "Storage: $STORAGE"
echo

# stop
docker stop auth-service
docker rm auth-service

# build
docker build -t rebel1l/auth-service:$BRANCH .

# start
docker run --name auth-service -d -it  -v $STORAGE:/usr/bin/app/storage -p 3000:3000 rebel1l/auth-service:$BRANCH
docker logs auth-service
docker ps
