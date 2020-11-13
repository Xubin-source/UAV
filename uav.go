/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	//"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Task struct {
	ID           string            `json:"ID"`
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
	Uri          string `json:"Uri"`
	EncryptedKey string `json:"EncryptedKey"`
	Tags         string `json:"Tags"`
}

//var status = [3]string{"pending", "processing", "done"}

//vote contracts

/*
votes中存放无人机ID，代表所投票的无人机
*/
//初始化账本
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	tasks := []Task{
		Task{ID: "task1", Commander: "A", NumOfDevices: "2", Status: "pending", TargetVotes: "2", Votes: map[string]string{"000": "stone", "001": "stone"}, Result: map[string]int{"stone": 2}},
		Task{ID: "task2", Commander: "B", NumOfDevices: "3", Status: "processing", TargetVotes: "2", Votes: map[string]string{"100": "stone", "101": "stone", "102": "flower"}, Result: map[string]int{"stone": 2, "flower": 0}},
	}

	for _, task := range tasks {
		taskJSON, _ := json.Marshal(task)

		err := ctx.GetStub().PutState(task.ID, taskJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err)
		}

	}
	return nil
}
func (s *SmartContract) CreateTask(ctx contractapi.TransactionContextInterface, taskId string, commander string, numOfDevices string, targetVotes string) error {
	task := Task{
		ID:           taskId,
		Commander:    commander,
		NumOfDevices: numOfDevices,
		TargetVotes:  targetVotes,
		Status:       "pending",
		Votes:        map[string]string{},
		Result:       map[string]int{},
	}
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(taskId, taskJSON)

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

	for _, result := range task.Votes { //遍历Map,统计结果
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
func (s *SmartContract) AddRecord(ctx contractapi.TransactionContextInterface, deviceId string, uri string, key string, tags string) error {
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
func (s *SmartContract) QueryData(ctx contractapi.TransactionContextInterface, uri string) (*Record, error) {
	recordAsBytes, err := ctx.GetStub().GetState(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if recordAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", uri)
	}

	record := new(Record)
	_ = json.Unmarshal(recordAsBytes, record)

	return record, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create uav chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting uav chaincode: %s", err.Error())
	}
}
