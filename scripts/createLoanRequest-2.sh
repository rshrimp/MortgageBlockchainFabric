#!/bin/bash

export FABRIC_START_WAIT=2

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
echo -e "----\e[5;32;40mNow 1st phase of the Loan Request to get fico, appraisal, title search and insurance quotes \e[m"
echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"


echo -e "................\e[5;32;40mInitiate Mortgage request on  lending ledger for the customer\e[m ........................"
docker exec cli.Bank bash -c "peer chaincode invoke -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"initiateMortgage\", \"TestCustomer891011\", \"891011\",\"1000000\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mAudit query lending ledger for the request created\e[m ........................"
docker exec cli.Audit bash -c "peer chaincode query -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"query\",\"TestCustomer891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mrequest ficoScore for  TestCustomer891011 \e[m ..................................."
docker exec cli.Bank bash -c "peer chaincode invoke -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"getFicoScores\", \"TestCustomer891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mAudit query lending ledger for the ficoScore created\e[m ..............................."
docker exec cli.Audit bash -c "peer chaincode query -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"query\",\"TestCustomer891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mrequest Appraisal \e[m........................"
docker exec cli.Appraiser bash -c "peer chaincode invoke -C books -n bookschaincode -v 0 -c '{\"Args\":[\"getAppraisal\", \"891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mAudit query books ledger for the appraiser created f\e[m..............................."
docker exec cli.Audit bash -c "peer chaincode query -C books -n bookschaincode -v 0 -c '{\"Args\":[\"query\",\"891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mrequest Insurance quote \e[m ..................................."
docker exec cli.Bank bash -c "peer chaincode invoke -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"getInsuranceQuote\", \"TestCustomer891011\", \"891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mAudit query lending ledger for the insurance quote created\e[m ..............................."
docker exec cli.Audit bash -c "peer chaincode query -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"query\",\"TestCustomer891011\"]}'"

echo -e "................\e[5;32;40mget Title on books ledger \e[m  .................................."
docker exec cli.Title bash -c "peer chaincode invoke -C books -n bookschaincode -v 0 -c '{\"Args\":[\"getTitle\", \"891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mAudit query books ledger for the Title created \e[m..............................."
docker exec cli.Audit bash -c "peer chaincode query -C books -n bookschaincode -v 0 -c '{\"Args\":[\"query\",\"891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
echo -e "----\e[5;32;40mNow 2nd phase of the Loan Request to close the loan change titles and update registry \e[m"
echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"

echo -e "................\e[5;32;40mclose the mortgage based on bank criteria \e[m..................................."
docker exec cli.Bank bash -c "peer chaincode invoke -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"closeMortgage\", \"TestCustomer891011\"]}'"
sleep ${FABRIC_START_WAIT}


echo -e "................\e[5;32;40mAudit query lending ledger for closed mortgage status \e[m..............................."
docker exec cli.Audit bash -c "peer chaincode query -C lending -n lendingchaincode -v 0 -c '{\"Args\":[\"query\",\"TestCustomer891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mChange Title on books ledger on \e[m  .................................."
docker exec cli.Title bash -c "peer chaincode invoke -C books -n bookschaincode -v 0 -c '{\"Args\":[\"changeTitle\", \"891011\", \"TestCustomer891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mAudit query books ledger for the Titlechanges \e[m..............................."
docker exec cli.Audit bash -c "peer chaincode query -C books -n bookschaincode -v 0 -c '{\"Args\":[\"query\",\"891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mrecord purchase on records ledger for new owner  \e[m ..................................."
docker exec cli.Registry bash -c "peer chaincode invoke -C records -n recordschaincode -v 0 -c '{\"Args\":[\"recordPurchase\", \"891011\"]}'"
sleep ${FABRIC_START_WAIT}

echo -e "................\e[5;32;40mAudit query records ledger for the owner changes if loan was successful \e[m..............................."
docker exec cli.Audit bash -c "peer chaincode query -C records -n recordschaincode -v 0 -c '{\"Args\":[\"query\",\"891011\"]}'"
sleep ${FABRIC_START_WAIT}
