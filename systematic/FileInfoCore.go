package systematic

import (
	"encoding/json"
	"errors"
	log "coldchian.go/lib/log4go"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func AddFile(stub shim.ChaincodeStubInterface, fileArr []string) (error, string) {
	//如果带有附件信息，就保存附件内容
	if len(fileArr) > 0 {
		for _, fileItem := range fileArr {
			var fileObj *FileInfo
			err := json.Unmarshal([]byte(fileItem), &fileObj)
			if err != nil {
				log.Error(ERRNMCC0001)
				return err, ERRNMCC0001
			}

			var fileId = fileObj.FileId
			//确保数据库中的id值保持唯一，校验唯一性
			tmpFileData, err := stub.GetState(fileId)
			if err != nil {
				log.Error(ERRNMCC0006)
				return err, ERRNMCC0006
			} else if tmpFileData != nil {
				var customError = errors.New(ERRNMCC0010)
				log.Error(ERRNMCC0010)
				return customError, ERRNMCC0010
			}

			//新增文件信息
			err = stub.PutState(fileId, []byte(fileItem))
			if err != nil {
				log.Error(ERRNMCC0005)
				return err, ERRNMCC0005
			}
		}
	}

	return nil, ""
}

func SelectFile(stub shim.ChaincodeStubInterface, fileId string) (error, string, string) {
	fileByteArr, err := stub.GetState(fileId)
	if err != nil {
		log.Error(ERRNMCC0006)
		return err, ERRNMCC0006, ""
	} else if fileByteArr == nil {
		var customError = errors.New(ERRNMCC0012)
		log.Error(ERRNMCC0012)
		return customError, ERRNMCC0012, ""
	}

	return nil, "", string(fileByteArr)
}
