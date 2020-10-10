// UAV
package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Task struct {
	Commander    string            `json:"Commander"`
	NumOfDevices string            `json:"NumOfDevices"`
	Status       []string          `json:"Status"`
	TargetVotes  string            `json:"TargetVotes"`
	Votes        map[string]string `json:"TargetVotes"` //无人机ID - 识别结果
	Result       map[string]int    `json:"Result"`      //识别结果 - 票数
}

type QueryResult struct { //查询任务
	Key    string `json:"Key"`
	Record *Task
}

type Record struct { //查询无人机的上传的数据
	DeviceId     string   `json:"DeviceId"`
	Uri          string   `json:"Uri"`
	EncryptedKey string   `json:"EncryptedKey"`
	Tags         []string `json:"Tags"`
}

var Status = [3]string{pending, processing, done}

//vote contracts

/*
votes中存放无人机ID，代表所投票的无人机
*/

func (s *SmartContract) CreateTask(ctx contractapi.TransactionContextInterface, taskId string, commander string, numOfDevices string, targetVotes string) error {
	task := Task{
		Commander:    commander,
		NumOfDevices: numOfDevices,
		TargetVotes:  targetVotes,
		Status:       Status[0],
		Votes:        {},
		Result:       {},
	}
	taskAsBytes, _ := json.Marshal(task)
	return ctx.GetStub().PutState(taskId, taskAsBytes)

}

//查询任务
func (s *SmartContract) QueryTask(ctx contractapi.TransactionContextInterface, taskId string) (*Task, error) {
	taskAsBytes, err := ctx.GetStub().GetState(taskId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if taskAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", taskId)
	}

	task := new(Task)
	_ = json.Unmarshal(taskAsBytes, task)

	return task, nil
}

//投票
func (s *SmartContract) Vote(ctx contractapi.TransactionContextInterface, taskId string, deviceId string, result string) error {
	task, err := s.QueryTask(ctx, taskId)

	if err != nil {
		return err
	}

	//task.Votes = append(task.Votes, deviceId)
	task.Votes[deviceId] = result

	taskAsBytes, _ := json.Marshal(task)

	return ctx.GetStub().PutState(taskId, taskAsBytes)
}

//计算投票结果
func (s *SmartContract) Calculate(ctx contractapi.TransactionContextInterface, taskId string) error {
	task, err := s.QueryTask(ctx, taskId)

	if err != nil {
		return err
	}

	for id, result := range task.Votes { //遍历Map,统计结果
		task.Result[result]++
	}

	taskAsBytes, _ := json.Marshal(task)

	return ctx.GetStub().PutState(taskId, taskAsBytes)

}

/*
func (s *SmartContract) CheckStatus(ctx contractapi.TransactionContextInterface) error {

}

func (s *SmartContract) CheckResult(ctx contractapi.TransactionContextInterface) error {

}
*/

//数据管理

//uri作为key
func (s *SmartContract) AddRecord(ctx contractapi.TransactionContextInterface, deviceId string, uri string, key string, tags []string) error {
	record := Record{
		DeviceId:     deviceId,
		Uri:          uri,
		EncryptedKey: key,
		Tags:         tags,
	}

	recordAsBytes, _ := json.Marshal(record)
	return ctx.GetStub().PutState(uri, recordAsBytes)

}

//查询数据，uri作为索引
func (s *SmartContract) QueryData(ctx contractapi.TransactionContextInterface, uri string) error {
	recordAsBytes, err := ctx.GetStub().GetState(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if recordAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", deviceId)
	}

	record := new(Record)
	_ = json.Unmarshal(recordAsBytes, record)

	return record, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
