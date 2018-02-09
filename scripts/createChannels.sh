#!/bin/bash

export FABRIC_START_WAIT=5

echo -e "-------------------------\e[5;32;40mNow Creating channels\e[m -----------------------------"
docker exec cli.Registry bash -c 'peer channel create -c records -f ./channels/Records.tx -o orderer.mrtgexchg.com:7050'
sleep ${FABRIC_START_WAIT}
docker exec cli.Registry bash -c 'peer channel join -b records.block'
sleep ${FABRIC_START_WAIT}

docker exec  cli.Bank bash -c 'peer channel create -c lending -f ./channels/Lending.tx -o orderer.mrtgexchg.com:7050'
sleep ${FABRIC_START_WAIT}
docker exec  cli.Bank bash -c 'peer channel join -b lending.block'
sleep ${FABRIC_START_WAIT}

docker exec  cli.Appraiser bash -c 'peer channel create -c books -f ./channels/Books.tx -o orderer.mrtgexchg.com:7050'
sleep ${FABRIC_START_WAIT}
docker exec  cli.Appraiser bash -c 'peer channel join -b books.block'

sleep ${FABRIC_START_WAIT}
echo -e "-------------------------\e[5;32;40mNow Joining channels\e[m -----------------------------"
#Registry  joins 2 channels, but we already joined records,  so join the other
docker exec cli.Registry bash -c 'peer channel join -b books.block'
sleep ${FABRIC_START_WAIT}

#bank  joins all channels, but we already joined lending when we created it,  so join the other two
docker exec cli.Bank bash -c 'peer channel join -b records.block'
sleep ${FABRIC_START_WAIT}
docker exec cli.Bank bash -c 'peer channel join -b books.block'
sleep ${FABRIC_START_WAIT}

#Appraiser  joins 2 channels, but we already joined books when we created it,  so join the other
docker exec cli.Appraiser bash -c 'peer channel join -b records.block'
sleep ${FABRIC_START_WAIT}

#Title  joins 3 channels
docker exec cli.Title bash -c 'peer channel join -b records.block'
sleep ${FABRIC_START_WAIT}
docker exec cli.Title bash -c 'peer channel join -b books.block'
sleep ${FABRIC_START_WAIT}
docker exec cli.Title bash -c 'peer channel join -b lending.block'
sleep ${FABRIC_START_WAIT}


#insurance  joins all channels
docker exec cli.Insurance bash -c 'peer channel join -b records.block'
sleep ${FABRIC_START_WAIT}
docker exec cli.Insurance bash -c 'peer channel join -b lending.block'
sleep ${FABRIC_START_WAIT}
docker exec cli.Insurance bash -c 'peer channel join -b books.block'
sleep ${FABRIC_START_WAIT}

#Audit  joins all channels
docker exec cli.Audit bash -c 'peer channel join -b records.block'
sleep ${FABRIC_START_WAIT}
docker exec cli.Audit bash -c 'peer channel join -b lending.block'
sleep ${FABRIC_START_WAIT}
docker exec cli.Audit bash -c 'peer channel join -b books.block'
sleep ${FABRIC_START_WAIT}

#Fico  joins 1 channels,
docker exec cli.Fico bash -c 'peer channel join -b lending.block'
sleep ${FABRIC_START_WAIT}


echo -e ".. \e[5;32;40mlet us use the anchor peer update transactions:\e[m"

docker exec cli.Bank bash -c 'peer channel update -o orderer.mrtgexchg.com:7050 -c lending -f ./channels/lendinganchor.tx'
sleep ${FABRIC_START_WAIT}
docker exec cli.Appraiser bash -c 'peer channel update -o orderer.mrtgexchg.com:7050 -c books -f ./channels/booksanchor.tx'
sleep ${FABRIC_START_WAIT}
docker exec cli.Registry bash -c 'peer channel update -o orderer.mrtgexchg.com:7050 -c records -f ./channels/recordsanchor.tx'
