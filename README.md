# UAV
测试：

查询账本：
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["QueryTask"，＂task1＂]}' 
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["QueryTask"，＂task2＂]}'

创建任务：
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["CreateTask","task3","X","3","2"]}'//指挥官X,设备数量３，目标票数２

投票：三个设备的识别结果
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["Vote"，＂task3＂,"000","flower"]}'//设备000，识别结果为花
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["Vote"，＂task3＂,"002","flower"]}'
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["Vote"，＂task3＂,"002","stone"]}'

计算：
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["Calculate"，＂task3＂]}'

添加数据：
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["AddRecord"，＂000＂,"uri","key","tags"]}'

查询数据：
peer chaincode query -C $CHANEL_NAME -n $CHAINCODE_NAME -c '{"Args":["QueryData"，＂uri＂]}'

