package model

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
	log "coldchain.go/lib/log4go"
	"coldchain.go/systematic"
)
//GENERATE_START
	func AddResourceBox(stub shim.ChaincodeStubInterface, ResourceBoxJson string, fileArr []string) (error, string){
		var ResourceBoxObj *ResourceBox
	
		err := json.Unmarshal([]byte(ResourceBoxJson), &ResourceBoxObj)
		if err != nil{
			log.Error(systematic.ERRNMCC0001)
			return err, systematic.ERRNMCC0001
		}
		
		var ResourceBoxKey = ResourceBoxObj.BoxId//获取主键id
	
		log.Info("Function Name:AddResourceBox ---- Message:PutState execute...")
		err = stub.PutState(ResourceBoxKey, []byte(ResourceBoxJson))
		if err != nil{
			log.Error(systematic.ERRNMCC0005)
			return err, systematic.ERRNMCC0005
		}else {
			log.Info("Function Name:AddResourceBox ---- Message:PutState SUCCESS!!!")
		}
		return nil, ""
	}
	func DelResourceBox(stub shim.ChaincodeStubInterface, ResourceBoxId string) (error, string) {
		log.Info("Function Name:DelResourceBox ---- Message:DelState execute...")
		//先取到数据，判断要删除的数据信息在数据库中存不存在
		tmpResourceBoxByteArr, err := stub.GetState(ResourceBoxId)
		if err != nil{
			log.Error(systematic.ERRNMCC0004)
			return err, systematic.ERRNMCC0004
		}else if tmpResourceBoxByteArr == nil{
			//说明没有对应的数据信息
			var customError = errors.New(systematic.ERRNMCC0012)
			log.Error(systematic.ERRNMCC0012)
			return customError, systematic.ERRNMCC0012
		}
	
		err = stub.DelState(ResourceBoxId)
		if err != nil{
			log.Error(systematic.ERRNMCC0004)
			return err, systematic.ERRNMCC0004
		}else {
			log.Info("Function Name:DelResourceBox ---- Message:DelState SUCCESS!!!")
		}
	
		return nil, ""
	}
	func UpdateResourceBox(stub shim.ChaincodeStubInterface, ResourceBoxJson string) (error, string) {
		var ResourceBoxObj *ResourceBox
	
		err := json.Unmarshal([]byte(ResourceBoxJson), &ResourceBoxObj)
		if err != nil{
			log.Error(systematic.ERRNMCC0001)
			return err, systematic.ERRNMCC0001
		}
		
		var ResourceBoxKey = ResourceBoxObj.BoxId//获取主键id
		tmpResourceBoxByteArr, err := stub.GetState(ResourceBoxKey)
		if err != nil{
			log.Error(systematic.ERRNMCC0006)
			return err, systematic.ERRNMCC0006
		}else if tmpResourceBoxByteArr == nil{
			var customError = errors.New(systematic.ERRNMCC0015)
			log.Error(systematic.ERRNMCC0015)
			return customError, systematic.ERRNMCC0015
		}
	
		log.Info("Function Name:UpdateResourceBox ---- Message:PutState execute...")
		err = stub.PutState(ResourceBoxKey, []byte(ResourceBoxJson))
		if err != nil{
			log.Error(systematic.ERRNMCC0011)
			return err, systematic.ERRNMCC0011
		}else {
			log.Info("Function Name:UpdateResourceBox ---- Message:PutState SUCCESS!!!")
		}
	
		return nil, ""
	}
	func SelectResourceBox(stub shim.ChaincodeStubInterface, ResourceBoxId string) (error, string, string) {
		ResourceBoxByteArr, err := stub.GetState(ResourceBoxId)
		if err != nil{
			log.Error(systematic.ERRNMCC0006)
			return err, systematic.ERRNMCC0006, ""
		}else if ResourceBoxByteArr == nil{
			var customError = errors.New(systematic.ERRNMCC0012)
			log.Error(systematic.ERRNMCC0012)
			return customError, systematic.ERRNMCC0012, ""
		}
	
		return nil, "", string(ResourceBoxByteArr)
	}
	func SelectAllResourceBox(stub shim.ChaincodeStubInterface) (error, string, []string) {
		var queryString = "{\"selector\": {\"dataType\": \"ResourceBox\"}}"
	
		resArr, err, errCode := systematic.ConditionQuery(stub, queryString)
		if err != nil{
			return  err, errCode, nil
		}
	
		return nil, "", resArr
	}
	
	func QueryResourceBox(stub shim.ChaincodeStubInterface, queryString string) (error, string, []string) {
		var pQueryString = fmt.Sprintf("{\"selector\": %s}", queryString)
	
		resArr, err, errCode := systematic.ConditionQuery(stub, pQueryString)
		if err != nil{
			return  err, errCode, nil
		}
	
		return nil, "", resArr
	}
//GENERATE_END