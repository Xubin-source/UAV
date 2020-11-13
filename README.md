# UAV
测试：
初始化账本：
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'
查询账本：
peer chaincode query -C mychannel -n basic -c '{"Args":["QueryTask"，＂taskid＂]}' //taskid类似于task1

创建任务：
peer chaincode query -C mychannel -n basic -c '{"Args":["CreateTask"，＂task３＂,"X","3","2"]}'//指挥官X,设备数量３，目标票数２

投票：三个设备的识别结果
peer chaincode query -C mychannel -n basic -c '{"Args":["Vote"，＂task3＂,"000","flower"]}'//设备000，识别结果为花
peer chaincode query -C mychannel -n basic -c '{"Args":["Vote"，＂task3＂,"002","flower"]}'
peer chaincode query -C mychannel -n basic -c '{"Args":["Vote"，＂task3＂,"002","stone"]}'

计算：
peer chaincode query -C mychannel -n basic -c '{"Args":["Calculate"，＂task3＂]}'

添加数据：
peer chaincode query -C mychannel -n basic -c '{"Args":["AddRecord"，＂000＂,"uri","key","tags"]}'

查询数据：
peer chaincode query -C mychannel -n basic -c '{"Args":["QueryData"，＂uri＂]}'

