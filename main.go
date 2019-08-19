package main

import (
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	log "coldchain.go/lib/log4go"
	"coldchain.go/appRouter"
	"coldchain.go/systematic"
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function =="Invoke" {
		return t.InvokeGate(stub, args)
	}else if function == "Query" {
		return t.QueryGate(stub, args)
	}else {
		log.Error("Entering errorGate name, please make sure!!!")
		return shim.Error("Entering errorGate name, please make sure!!!")
	}

	return shim.Success(nil)
}

//数据操作入口
func (t *SimpleChaincode)InvokeGate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	msg, err := getMessageFromArgs(args)
	if err != nil{
		returnErrMsg, _ := json.Marshal(msg)
		return shim.Success(returnErrMsg)
	}

	//初始化返回报文信息
	returnMsg := systematic.InitMessage(msg.Id, msg.Channel, msg.TranCode, msg.TranDate, msg)
	appRouter.ProcessBusiness_Invoke(stub, msg.TranCode, msg, &returnMsg)

	returnMsgByte, _ := json.Marshal(returnMsg)
	return shim.Success(returnMsgByte)
}

//数据查询入口
func (t *SimpleChaincode)QueryGate(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	msg, err := getMessageFromArgs(args)
	if err != nil{
		returnErrMsg, _ := json.Marshal(msg)
		return shim.Success(returnErrMsg)
	}

	returnMsg := systematic.InitMessage(msg.Id, msg.Channel, msg.TranCode, msg.TranDate, msg)
	appRouter.ProcessBusiness_Query(stub, msg.TranCode, msg, &returnMsg)

	returnMsgByte, _ := json.Marshal(returnMsg)
	return shim.Success(returnMsgByte)
}

func getMessageFromArgs(args []string) (*systematic.Message, error) {
	var msg *systematic.Message
	if len(args) == 0{
		log.Error(systematic.ERRNMCC0014)
		var customErr = errors.New("请求报文args属性为空")

		var returnMsg systematic.Message
		returnMsg.Id = "Error args"
		returnMsg.Channel = "null"
		returnMsg.TranCode = "null"
		returnMsg.TranDate = systematic.GetTimeStamp()
		returnMsg.Debug = "Args hasn't data information"

		returnMsg.Type = "1" //0是请求报文，1是返回报文
		returnMsg.RetCode = "11" //11返回的是错误报文
		return &returnMsg, customErr
	}else {
		msgStr := args[0]
		err := json.Unmarshal([]byte(msgStr), &msg)
		if err != nil{
			var returnMsg systematic.Message
			returnMsg.Id = "Error msg"
			returnMsg.Channel = "null"
			returnMsg.TranCode = "null"
			returnMsg.TranDate = systematic.GetTimeStamp()
			returnMsg.Debug = "Unmarshal 'message' obj has error"

			returnMsg.Type = "1" //0是请求报文，1是返回报文
			returnMsg.RetCode = "11" //11返回的是错误报文
			return &returnMsg, err
		}
		return msg, nil
	}
	return msg, nil
}

func main() {
	//log4go工具使用
	systematic.SetConsoleAndFileLog()
	defer log.Close()

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting car chaincode: %s", err)
	}
}