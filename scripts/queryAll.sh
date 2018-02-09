#!/bin/bash

export FABRIC_START_WAIT=20

docker exec cli.Audit bash -c "peer chaincode query -C records -n mrtgexchg -v 0 -c '{\"Args\":[\"queryAll\"]}'"
sleep ${FABRIC_START_WAIT}
docker exec cli.Audit bash -c "peer chaincode query -C lending -n mrtgexchg -v 0 -c '{\"Args\":[\"queryAll\"]}'"
sleep ${FABRIC_START_WAIT}
docker exec cli.Audit bash -c "peer chaincode query -C books -n mrtgexchg -v 0 -c '{\"Args\":[\"queryAll\"]}'"
