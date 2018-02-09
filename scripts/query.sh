#!/bin/bash

export FABRIC_START_WAIT=20

docker exec cli.Audit bash -c "peer chaincode query -C records -n mrtgexchg -v 0 -c '{\"Args\":[\"query\",\"456\"]}'"
sleep ${FABRIC_START_WAIT}
docker exec cli.Audit bash -c "peer chaincode query -C lending -n mrtgexchg -v 0 -c '{\"Args\":[\"query\",\"123\"]}'"
sleep ${FABRIC_START_WAIT}
docker exec cli.Audit bash -c "peer chaincode query -C books -n mrtgexchg -v 0 -c '{\"Args\":[\"query\",\"123\"]}'"
