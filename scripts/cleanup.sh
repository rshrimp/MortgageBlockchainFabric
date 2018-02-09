#!/bin/bash

export SLEEP_TIME=2


docker kill $(docker ps -q)
sleep {SLEEP_TIME}
docker rm $(docker ps -aq)
sleep {SLEEP_TIME}
docker rmi $(docker images dev-* -q)
