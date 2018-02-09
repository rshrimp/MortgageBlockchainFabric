
#!/bin/bash

export FABRIC_START_WAIT=5

echo -e '-----------------------\e[5;32;40m Install chaincodes\e[m---------------------------------------------------------'

echo " ----------------------------- For Records channel --------------------------------------------"
docker exec cli.Registry bash -c 'peer chaincode install -p mrtgexchg -n recordschaincode -v 0'
docker exec cli.Audit bash -c 'peer chaincode install -p mrtgexchg -n recordschaincode -v 0'
docker exec cli.Bank bash -c 'peer chaincode install -p mrtgexchg -n recordschaincode -v 0'
docker exec cli.Appraiser bash -c 'peer chaincode install -p mrtgexchg -n recordschaincode -v 0'
docker exec cli.Title bash -c 'peer chaincode install -p mrtgexchg -n recordschaincode -v 0'
docker exec cli.Insurance bash -c 'peer chaincode install -p mrtgexchg -n recordschaincode -v 0'

sleep ${FABRIC_START_WAIT}

echo "----------------------------- For books channel------------------------------------------------"
docker exec cli.Appraiser bash -c 'peer chaincode install -p mrtgexchg -n bookschaincode -v 0'
docker exec cli.Title bash -c 'peer chaincode install -p mrtgexchg -n bookschaincode -v 0'
docker exec cli.Bank bash -c 'peer chaincode install -p mrtgexchg -n bookschaincode -v 0'
docker exec cli.Insurance bash -c 'peer chaincode install -p mrtgexchg -n bookschaincode -v 0'
docker exec cli.Audit bash -c 'peer chaincode install -p mrtgexchg -n bookschaincode -v 0'
docker exec cli.Registry bash -c 'peer chaincode install -p mrtgexchg -n bookschaincode -v 0'
sleep ${FABRIC_START_WAIT}

echo " ----------------- For lending channel ----------------------------------------------------"
docker exec cli.Bank bash -c 'peer chaincode install -p mrtgexchg -n lendingchaincode -v 0'
docker exec cli.Fico bash -c 'peer chaincode install -p mrtgexchg -n lendingchaincode -v 0'
docker exec cli.Insurance bash -c 'peer chaincode install -p mrtgexchg -n lendingchaincode -v 0'
docker exec cli.Audit bash -c 'peer chaincode install -p mrtgexchg -n lendingchaincode -v 0'
docker exec cli.Title bash -c 'peer chaincode install -p mrtgexchg -n lendingchaincode -v 0'
sleep ${FABRIC_START_WAIT}

echo -e "-----------------------'\e[5;32;40m Instantiate chaincodes\e[m---------------------------------------------------------"

echo "---------------Instantiate chaincode on lending channel with permission for Bank, Insurance and Fico to invoke the chaincode and Audit+Title to have readonly access-------------------"

docker exec cli.Bank bash -c "peer chaincode instantiate -C lending -n lendingchaincode -v 0 -c '{\"Args\":[]}'  -P \"OR ('BankMSP.member', 'InsuranceMSP.member','FicoMSP.member')\""
echo "---------------Instantiate chaincode on books channel with permission for Appraiser&  Title  to invoke the chaincode and Bank+Insurance+Audit+Registry to have readonly access----------------------------"
docker exec cli.Appraiser bash -c "peer chaincode instantiate -C books -n bookschaincode -v 0 -c '{\"Args\":[]}' -P \"OR ('AppraiserMSP.member', 'TitleMSP.member')\""
echo "---------------Instantiate chaincode on records channel with permission for Registry only  to invoke the chaincode and all others to have readonly access----------------------------"
docker exec cli.Registry bash -c "peer chaincode instantiate -C records -n recordschaincode -v 0 -c '{\"Args\":[]}'"

echo -e "----------------------'\e[5;32;40m END\e[m\'---------------------------------------------------------"
