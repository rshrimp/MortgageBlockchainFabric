#!/bin/bash

export FABRIC_START_WAIT=5
echo -e '******\e[5;32;40mList of orgs and their subscriptions to lending books and records blockchain\e[m************************************'

echo -e '******\e[5;32;40mBank (lending, records, books channels)\e[m************************************'
docker exec cli.Bank bash -c 'peer channel list'
sleep ${FABRIC_START_WAIT}
echo -e '******\e[5;32;40mInsurance  (lending, records, books channels)\e[m************************************'
docker exec cli.Insurance bash -c 'peer channel list'
sleep ${FABRIC_START_WAIT}
echo -e '******\e[5;32;40mRegistry ( records channel)\e[m************************************'
docker exec cli.Registry bash -c 'peer channel list'
sleep ${FABRIC_START_WAIT}
echo -e '******\e[5;32;40mTitle ( records, books channels)\e[m************************************'
docker exec cli.Title bash -c 'peer channel list'
sleep ${FABRIC_START_WAIT}
echo -e '******\e[5;32;40mFico (lending, channel)\e[m************************************'
docker exec cli.Fico bash -c 'peer channel list'
sleep ${FABRIC_START_WAIT}
echo -e '******\e[5;32;40mAppraiser (records, books channels)\e[m************************************'
docker exec cli.Appraiser bash -c 'peer channel list'
sleep ${FABRIC_START_WAIT}
echo -e '******\e[5;32;40mAudit (lending, records, books channels)\e[m************************************'
docker exec cli.Audit bash -c 'peer channel list'
