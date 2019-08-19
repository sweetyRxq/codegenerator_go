package appRouter

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	log "coldchain.go/lib/log4go"
	"coldchain.go/model"
	"coldchain.go/systematic"
)
//GENERATE_START

func ProcessBusiness_Invoke(stub shim.ChaincodeStubInterface, tranCode string, requestMsg *systematic.Message, returnMsg *systematic.Message ) {

	switch tranCode{
	case "CreateResourceCar":
		log.Info("Request Message Id:" + requestMsg.Id + " CreateResourceCar perform start")

		//判断新增数据是否涉及到了附件的上传操作
		var hasFile string
		for _, argsItem := range requestMsg.Args{
			if argsItem.Name == "hasFile"{
				hasFile = argsItem.Value
			}
		}

		var resourceCarJSON string
		var fileArr []string
		//如果涉及到了附件信息
		if hasFile == "true"{
			dataArr, err := systematic.GetDataFromMessage(requestMsg) //获取到的是个data数组
			if err != nil{
				systematic.CombinationErrorObj(err.Error(), "CreateResourceCar", systematic.ERRNMCC0012, returnMsg)
				systematic.ErrorMessage(returnMsg)
			}

			for _, dataItem := range dataArr{
				if dataItem.DataType == "ResourceCar"{
					resourceCarJSON = dataItem.Content
				}else if dataItem.DataType == "FileInfo"{
					fileArr = append(fileArr, dataItem.Content)
				}
			}
		}else if hasFile == "false"{
			//如果是单纯的新增业务数据
			tmpResourceCarJSON, err := systematic.GetOneDataFromMessage(requestMsg)
			if err != nil{
				systematic.CombinationErrorObj(err.Error(), "CreateResourceCar", systematic.ERRNMCC0012, returnMsg)
				systematic.ErrorMessage(returnMsg)
			}

			resourceCarJSON = tmpResourceCarJSON
		}

		err, errCode := model.AddResourceCar(stub, resourceCarJSON, fileArr)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "CreateResourceCar", errCode, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}

		log.Info("Return Message Id:" + returnMsg.Id + " CreateResourceCar perform finish")
		break
	case "DelResourceCar":
		log.Info("Request Message Id:" + requestMsg.Id + " DelResourceCar perform start")
		var args = systematic.GetArgsOfDelOrQuery(requestMsg)
		var id string

		for _, argItem := range args{
			switch argItem.Name {
			case "id":
				id = argItem.Value
				break
			}
		}

		err, errCode := model.DelResourceCar(stub, id)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "DelResourceCar", errCode, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}
		log.Info("Return Message Id:" + returnMsg.Id + " DelResourceCar perform finish")
		break
	case "UpdateResourceCar":
		log.Info("Request Message Id:" + requestMsg.Id + " UpdateResourceCar perform start")
		resourceCarJSON, err := systematic.GetOneDataFromMessage(requestMsg)
		if err != nil{
			log.Error(systematic.ERRNMCC0012)
			systematic.CombinationErrorObj(err.Error(), "UpdateResourceCar", systematic.ERRNMCC0012, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}

		err, errCode := model.UpdateResourceCar(stub, resourceCarJSON)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "UpdateResourceCar", errCode, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}
		log.Info("Return Message Id:" + returnMsg.Id + " UpdateResourceCar perform finish")
		break
	case "FileInvoke":
		log.Info("Request Message Id:" + requestMsg.Id + " FileInvoke perform start")

		var fileArr []string

		dataArr, err := systematic.GetDataFromMessage(requestMsg) //获取到的是个data数组
		if err != nil {
			systematic.CombinationErrorObj(err.Error(), "FileInvoke", systematic.ERRNMCC0012, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}

		for _, dataItem := range dataArr {
			if dataItem.DataType == "FileInfo" {
				fileArr = append(fileArr, dataItem.Content)
			}
		}

		err, errCode := systematic.AddFile(stub, fileArr)
		if err != nil {
			systematic.CombinationErrorObj(err.Error(), "FileInvoke", errCode, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}

		log.Info("Return Message Id:" + returnMsg.Id + " FileInvoke perform finish")
		break
	default:
		var errDescription = fmt.Sprintf("Incorrect invoke transaction type routing function [ %s ]!", tranCode)
		log.Error(errDescription)
		systematic.CombinationErrorObj(errDescription, "ProcessBusiness_Invoke", "ERR_METHOD_NOT_FOUND", returnMsg)
		systematic.ErrorMessage(returnMsg)
	}
}

func ProcessBusiness_Query(stub shim.ChaincodeStubInterface, tranCode string, requestMsg *systematic.Message, returnMsg *systematic.Message ) {
	switch tranCode{
	case "GetVersion": //获取到版本号的接口
		log.Info("Request Message Id:" + requestMsg.Id + " GetVersion perform start")
		systematic.GetVersion(stub, requestMsg, returnMsg)
		log.Info("Return Message Id:" + returnMsg.Id + " GetVersion perform finish")
		break
	case "SelectAllResourceCar":
		log.Info("Request Message Id:" + requestMsg.Id + " SelectAllResourceCar perform start")
		err, errCode, resArr := model.SelectAllResourceCar(stub)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "SelectAllResourceCar", errCode, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}else {
			for _, item := range resArr{
				var data systematic.Data
				data.Content = item
				data.DataType = "ResourceCar"

				returnMsg.Data = append(returnMsg.Data, data)
			}
		}
		log.Info("Return Message Id:" + returnMsg.Id + " SelectAllResourceCar perform finish")
		break
	case "SelectResourceCar":
		log.Info("Request Message Id:" + requestMsg.Id + " SelectResourceCar perform start")
		resourceCarJSON, err := systematic.GetOneDataFromMessage(requestMsg)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "SelectResourceCar", systematic.ERRNMCC0012, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}

		var resourceCarObj *model.ResourceCar
		err = json.Unmarshal([]byte(resourceCarJSON), &resourceCarObj)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "SelectResourceCar", systematic.ERRNMCC0001, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}else {
			
			var resourceCarKey = resourceCarObj.CarId//获取主键id
			if resourceCarKey == ""{
				systematic.CombinationErrorObj(err.Error(), "SelectResourceCar", systematic.ERRNMCC0012, returnMsg)
				systematic.ErrorMessage(returnMsg)
			}else {
				err, errCode, result := model.SelectResourceCar(stub, resourceCarKey)
				if err != nil{
					systematic.CombinationErrorObj(err.Error(), "SelectResourceCar", errCode, returnMsg)
					systematic.ErrorMessage(returnMsg)
				}

				var data systematic.Data
				data.DataType = "ResourceCar"
				data.Content = result

				returnMsg.Data = append(returnMsg.Data, data)
			}
		}

		log.Info("Return Message Id:" + returnMsg.Id + " SelectResourceCar perform finish")
		break
	case "QueryResourceCar":
		log.Info("Request Message Id:" + requestMsg.Id + " QueryResourceCar perform start")
		// 获取QueryString
		var args = systematic.GetArgsOfDelOrQuery(requestMsg)
		var queryString string
		for _, argsItem := range args {
			switch argsItem.Name {
			case "queryString":
				queryString = argsItem.Value
				break
			}
		}
		if queryString == "" {
			tempQueryString, err := systematic.GetOneDataFromMessage(requestMsg)
			if err != nil {
				queryString = tempQueryString
			}
		}
		// 如果queryString为空
		if queryString == "" {
			queryString = "{\"dataType\":\"ResourceCar\"}"
		}
		err, errCode, resArr := model.QueryResourceCar(stub, queryString)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "SelectAllResourceCar", errCode, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}else {
			for _, item := range resArr{
				var data systematic.Data
				data.Content = item
				data.DataType = "ResourceCar"

				returnMsg.Data = append(returnMsg.Data, data)
			}
		}
		log.Info("Return Message Id:" + returnMsg.Id + " QueryResourceCar perform finish")
		break
	case "PaginateResourceCar":
		log.Info("Request Message Id:" + requestMsg.Id + " PaginateResourceCar perform start")
		systematic.PaginateQuery(stub, requestMsg, returnMsg)
		log.Info("Return Message Id:" + returnMsg.Id + " PaginateResourceCar perform finish")
		break;
	case "FileQuery":
		log.Info("Request Message Id:" + requestMsg.Id + " FileQuery perform start")
		fileJSON, err := systematic.GetOneDataFromMessage(requestMsg)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "FileQuery", systematic.ERRNMCC0012, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}

		var fileObj *systematic.FileInfo
		err = json.Unmarshal([]byte(fileJSON), &fileObj)
		if err != nil{
			systematic.CombinationErrorObj(err.Error(), "FileQuery", systematic.ERRNMCC0001, returnMsg)
			systematic.ErrorMessage(returnMsg)
		}else {
			var fileKey = fileObj.FileId
			if fileKey == ""{
				systematic.CombinationErrorObj(err.Error(), "FileQuery", systematic.ERRNMCC0012, returnMsg)
				systematic.ErrorMessage(returnMsg)
			}else {
				err, errCode, result := systematic.SelectFile(stub, fileKey)
				if err != nil{
					systematic.CombinationErrorObj(err.Error(), "FileQuery", errCode, returnMsg)
					systematic.ErrorMessage(returnMsg)
				}

				if result != ""{
					var data systematic.Data
					data.DataType = "FileInfo"
					data.Content = result

					returnMsg.Data = append(returnMsg.Data, data)
				}
			}
		}

		log.Info("Return Message Id:" + returnMsg.Id + " FileQuery perform finish")
		break
	default:
		var errDescription = fmt.Sprintf("Incorrect query transaction type routing function [ %s ]!", tranCode)
		log.Error(errDescription)
		systematic.CombinationErrorObj(errDescription, "ProcessBusiness_Query", "ERR_METHOD_NOT_FOUND", returnMsg)
		systematic.ErrorMessage(returnMsg)
	}
}

//GENERATE_END
