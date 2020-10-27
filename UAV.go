package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type SimpleChaincode struct {
}

type Task struct {
	Commander    string            `json:"Commander"`
	NumOfDevices string            `json:"NumOfDevices"`
	Status       string            `json:"Status"`
	TargetVotes  string            `json:"TargetVotes"`
	Votes        map[string]string `json:"Votes"`  //无人机ID - 识别结果
	Result       map[string]int    `json:"Result"` //识别结果 - 票数
}

type QueryResult struct { //查询任务
	Key    string `json:"Key"`
	Record *Task
}

type Record struct { //查询无人机的上传的数据
	DeviceId     string `json:"DeviceId"`
	EncryptedKey string `json:"EncryptedKey"`
	Tags         string `json:"Tags"`
}

//var status = [3]string{"pending", "processing", "done"}

//vote contracts

/*
votes中存放无人机ID，代表所投票的无人机
*/

func (t *SimpleChaincode) Init(APIstub shim.ChaincodeStubInterface) pb.Response {

	task := []Task{
		Task{Commander: "A", NumOfDevices: "2", Status: "pending", TargetVotes: "2", Votes: map[string]string{"000": "stone", "001": "stone"}, Result: map[string]int{"stone": 2}},
		Task{Commander: "B", NumOfDevices: "3", Status: "processing", TargetVotes: "2", Votes: map[string]string{"100": "stone", "101": "stone", "102": "flower"}, Result: map[string]int{"stone": 2, "flower": 0}},
	}

	for i, tas := range task {
		taskinfoBytes, _ := json.Marshal(tas)
		APIstub.PutState("task"+strconv.Itoa(i), taskinfoBytes)
	}

	return shim.Success(nil)

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	if os.Getenv("DEVMODE_ENABLED") != "" {
		fmt.Println("invoking in devmode")
	}

	fun, args := stub.GetFunctionAndParameters()

	if fun == "CreateTask" {
		return t.CreateTask(stub, args)
	} else if fun == "QueryTask" {
		return t.QueryTask(stub, args)
	} else if fun == "Vote" {
		return t.Vote(stub, args)
	} else if fun == "Calculate" {
		return t.Calculate(stub, args)
	} else if fun == "QueryData" {
		return t.QueryData(stub, args)
	} else if fun == "AddRecord" {
		return t.AddRecord(stub, args)
	}
	return shim.Error("Recevied unkown function invocation")
}

//通过taskid创建
func (t *SimpleChaincode) CreateTask(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Please input correct information of person")
	}
	task := Task{
		Commander:    args[1],
		NumOfDevices: args[2],
		TargetVotes:  args[3],
		Status:       args[4],
		Votes:        map[string]string{},
		Result:       map[string]int{},
	}
	taskAsBytes, err := json.Marshal(task)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], taskAsBytes) //arg[0]是taksid
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(taskAsBytes)

}

//查询任务 通过taskid查询
func (t *SimpleChaincode) QueryTask(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Please input the correct number ")
	}

	task, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(task)

}

//投票
func (t *SimpleChaincode) Vote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 { //taskid,deviceid,result
		return shim.Error("Please input the correct number ")
	}

	var taskinfo Task
	//根据taskid获得task信息并转换为结构体格式taskinfo
	task, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed tqo get state")
	}
	_ = json.Unmarshal(task, &taskinfo)

	//获取识别结果
	taskinfo.Votes[args[1]] = args[2]

	//写入账本
	taskAsBytes, _ := json.Marshal(taskinfo)
	_ = stub.PutState(args[0], taskAsBytes)

	return shim.Success(taskAsBytes)
}

//计算投票结果
func (t *SimpleChaincode) Calculate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 { //taskid
		return shim.Error("Please input the correct number ")
	}
	var taskinfo Task

	task, err := stub.GetState(args[0]) //task是json格式
	if err != nil {
		return shim.Error("Failed to get state")
	}
	_ = json.Unmarshal(task, &taskinfo)

	for _, result := range taskinfo.Votes { //遍历Map,统计结果  result是votes中的识别结果

		taskinfo.Result[result]++
	}

	taskAsBytes, _ := json.Marshal(taskinfo)
	_ = stub.PutState(args[0], taskAsBytes)

	return shim.Success(taskAsBytes)
}

/*
func (s *SmartContract) CheckStatus(ctx contractapi.TransactionContextInterface) error {
}
func (s *SmartContract) CheckResult(ctx contractapi.TransactionContextInterface) error {
}
*/

//数据管理

//uri作为key
func (t *SimpleChaincode) AddRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 { //
		return shim.Error("Please input the correct number ")
	}

	record := Record{
		DeviceId:     args[1],
		EncryptedKey: args[2],
		Tags:         args[3],
	}

	recordAsBytes, _ := json.Marshal(record)
	_ = stub.PutState(args[0], recordAsBytes)

	return shim.Success(nil)

}

//查询数据，uri作为索引
func (t *SimpleChaincode) QueryData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Please input the correct number ")
	}

	record, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(record)

}

/*
func main() {

	err := shim.Start(new(UavChaincode))
	if err != nil {
		fmt.Printf("Error starting uav chaincode: %s ", err)
	}
}
*/
