#!/bin/bash
SCRIPTS="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# start redis
$SCRIPTS/restartRedis.sh

# start service
$SCRIPTS/restartService.sh
