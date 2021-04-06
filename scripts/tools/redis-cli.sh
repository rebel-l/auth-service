#!/bin/bash
 docker run -it --name redis-cli --link redis:redis --rm redis redis-cli -h redis -p 6379
