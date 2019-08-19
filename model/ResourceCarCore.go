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
	func AddResourceCar(stub shim.ChaincodeStubInterface, ResourceCarJson string, fileArr []string) (error, string){
		var ResourceCarObj *ResourceCar
	
		err := json.Unmarshal([]byte(ResourceCarJson), &ResourceCarObj)
		if err != nil{
			log.Error(systematic.ERRNMCC0001)
			return err, systematic.ERRNMCC0001
		}
		
		var ResourceCarKey = ResourceCarObj.CarId//获取主键id
	
		log.Info("Function Name:AddResourceCar ---- Message:PutState execute...")
		err = stub.PutState(ResourceCarKey, []byte(ResourceCarJson))
		if err != nil{
			log.Error(systematic.ERRNMCC0005)
			return err, systematic.ERRNMCC0005
		}else {
			log.Info("Function Name:AddResourceCar ---- Message:PutState SUCCESS!!!")
		}
		return nil, ""
	}
	func DelResourceCar(stub shim.ChaincodeStubInterface, ResourceCarId string) (error, string) {
		log.Info("Function Name:DelResourceCar ---- Message:DelState execute...")
		//先取到数据，判断要删除的数据信息在数据库中存不存在
		tmpResourceCarByteArr, err := stub.GetState(ResourceCarId)
		if err != nil{
			log.Error(systematic.ERRNMCC0004)
			return err, systematic.ERRNMCC0004
		}else if tmpResourceCarByteArr == nil{
			//说明没有对应的数据信息
			var customError = errors.New(systematic.ERRNMCC0012)
			log.Error(systematic.ERRNMCC0012)
			return customError, systematic.ERRNMCC0012
		}
	
		err = stub.DelState(ResourceCarId)
		if err != nil{
			log.Error(systematic.ERRNMCC0004)
			return err, systematic.ERRNMCC0004
		}else {
			log.Info("Function Name:DelResourceCar ---- Message:DelState SUCCESS!!!")
		}
	
		return nil, ""
	}
	func UpdateResourceCar(stub shim.ChaincodeStubInterface, ResourceCarJson string) (error, string) {
		var ResourceCarObj *ResourceCar
	
		err := json.Unmarshal([]byte(ResourceCarJson), &ResourceCarObj)
		if err != nil{
			log.Error(systematic.ERRNMCC0001)
			return err, systematic.ERRNMCC0001
		}
		
		var ResourceCarKey = ResourceCarObj.CarId//获取主键id
		tmpResourceCarByteArr, err := stub.GetState(ResourceCarKey)
		if err != nil{
			log.Error(systematic.ERRNMCC0006)
			return err, systematic.ERRNMCC0006
		}else if tmpResourceCarByteArr == nil{
			var customError = errors.New(systematic.ERRNMCC0015)
			log.Error(systematic.ERRNMCC0015)
			return customError, systematic.ERRNMCC0015
		}
	
		log.Info("Function Name:UpdateResourceCar ---- Message:PutState execute...")
		err = stub.PutState(ResourceCarKey, []byte(ResourceCarJson))
		if err != nil{
			log.Error(systematic.ERRNMCC0011)
			return err, systematic.ERRNMCC0011
		}else {
			log.Info("Function Name:UpdateResourceCar ---- Message:PutState SUCCESS!!!")
		}
	
		return nil, ""
	}
	func SelectResourceCar(stub shim.ChaincodeStubInterface, ResourceCarId string) (error, string, string) {
		ResourceCarByteArr, err := stub.GetState(ResourceCarId)
		if err != nil{
			log.Error(systematic.ERRNMCC0006)
			return err, systematic.ERRNMCC0006, ""
		}else if ResourceCarByteArr == nil{
			var customError = errors.New(systematic.ERRNMCC0012)
			log.Error(systematic.ERRNMCC0012)
			return customError, systematic.ERRNMCC0012, ""
		}
	
		return nil, "", string(ResourceCarByteArr)
	}
	func SelectAllResourceCar(stub shim.ChaincodeStubInterface) (error, string, []string) {
		var queryString = "{\"selector\": {\"dataType\": \"ResourceCar\"}}"
	
		resArr, err, errCode := systematic.ConditionQuery(stub, queryString)
		if err != nil{
			return  err, errCode, nil
		}
	
		return nil, "", resArr
	}
	
	func QueryResourceCar(stub shim.ChaincodeStubInterface, queryString string) (error, string, []string) {
		var pQueryString = fmt.Sprintf("{\"selector\": %s}", queryString)
	
		resArr, err, errCode := systematic.ConditionQuery(stub, pQueryString)
		if err != nil{
			return  err, errCode, nil
		}
	
		return nil, "", resArr
	}
//GENERATE_END