# UAV
测试：初始化账本中有两个任务，分别为task1和task2．

注意把$CHANEL_NAME 和$CHAINCODE_NAME改成对应的通道名称和链码名称

查询账本：
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["QueryTask","task1"]}' 

peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["QueryTask","task2"]}' 


创建任务：

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANEL_NAME -n $CHAINCODE_NAME --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"CreateTask","Args":["task3","x","3","2"]}'    //指挥官X,设备数量３，目标票数２

投票：三个设备的识别结果．分别用设备001,002,003对识别结果进行投票
//设备001投票
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANEL_NAME  -n $CHAINCODE_NAME --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Vote","Args":["task3","001","flower"]}'      //设备00１，识别结果为花，其余的只需要改第二个和第三个参数

计算：将上述三个设备的投票结果进行计算,存到result中

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANEL_NAME  -n $CHAINCODE_NAME --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Calculate","Args":["task3"]}'

添加数据：

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANEL_NAME  -n $CHAINCODE_NAME --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"AddRecord","Args":["000","uri","key","tags"]}'　 

查询数据：

peer chaincode query -C $CHANEL_NAME  -n $CHAINCODE_NAME -c '{"Args":["QueryData","uri"]}'

