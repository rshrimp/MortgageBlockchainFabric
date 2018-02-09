#!/bin/bash

export FABRIC_START_WAIT=3
export FABRIC_CFG_PATH=./

echo -e "\e[5;32;40mgenerating certificates in crypto-config folder for all entities\e[m "
cryptogen generate --config crypto-config.yaml
sleep ${FABRIC_START_WAIT}



echo -e "\e[5;32;40mgenerating geneis block\e[m "
mkdir orderer
configtxgen -profile MRTGEXCHGOrdererGenesis -outputBlock ./orderer/genesis.block

echo -e "\e[5;32;40mcreate the channel configuration blocks with this configuration file, by using the other profiles\e[m "

mkdir channels
configtxgen -profile RecordsChannel -outputCreateChannelTx ./channels/Records.tx -channelID records
sleep ${FABRIC_START_WAIT}
configtxgen -profile LendingChannel -outputCreateChannelTx ./channels/Lending.tx -channelID lending
sleep ${FABRIC_START_WAIT}
configtxgen -profile BooksChannel -outputCreateChannelTx ./channels/Books.tx -channelID books
sleep ${FABRIC_START_WAIT}

echo -e "\e[5;32;40mgenerate the anchor peer update transactions\e[m "


configtxgen -profile RecordsChannel -outputAnchorPeersUpdate ./channels/recordsanchor.tx -channelID records -asOrg RegistryMSP
sleep ${FABRIC_START_WAIT}
configtxgen -profile LendingChannel -outputAnchorPeersUpdate ./channels/lendinganchor.tx -channelID lending -asOrg BankMSP
sleep ${FABRIC_START_WAIT}
configtxgen -profile BooksChannel -outputAnchorPeersUpdate ./channels/booksanchor.tx -channelID books -asOrg AppraiserMSP
