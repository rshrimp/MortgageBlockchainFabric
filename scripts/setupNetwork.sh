#!/bin/bash


./scripts/createChannels.sh
sleep 5

./scripts/checkOrgChannelSubscription.sh
sleep 5

./scripts/chaincodeInstallInstantiate.sh
sleep 5

./scripts/createLedgerEntries.sh
sleep 5
