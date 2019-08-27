package systematic

import(
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
	"strings"
	"fmt"
	"bytes"
	"errors"
	"strconv"
	"encoding/json"
	log "coldchain.go/lib/log4go"
)

/*
	log信息显示的优先级
	FATAL //1
	ERROR //2
	WARNING //3
	INFO  //4
	DEBUG //5
*/

//终端和文件都写log信息
func SetConsoleAndFileLog()  {
	w1 := log.NewFileWriter()
	w1.SetPathPattern("/tmp/ChainCodeLog/Log%Y%M%D%H.log")
	//w1.SetPathPattern("/tmp/ChainCodeInfo.log")

	w2 := log.NewConsoleWriter()

	log.Register(w1)
	log.Register(w2)
	log.SetLevel(log.DEBUG)
}

//用户查询chainCode版本信息的接口
func GetVersion(stub shim.ChaincodeStubInterface, requestMsg *Message, returnMsg *Message) {
	var data Data
	data.DataType = "version"
	data.Content = "V1.0.0"

	returnMsg.Data = append(returnMsg.Data, data)
	return
}

/*
	log信息显示的优先级
	FATAL //1
	ERROR //2
	WARNING //3
	INFO  //4
	DEBUG //5
*/

//错误执行流程
//内部逻辑：该方法内会把错误信息，起始报错的方法名组成一个Error对象插入到返回报文的Error对象数组里面
func CombinationErrorObj(errText string, functionName string, errorCode string, returnMsg *Message) {
	var err Error
	err.Statement = errText
	err.Function = functionName
	err.Code = errorCode

	returnMsg.Errors = append(returnMsg.Errors, err)
}

/*
	Args   []Args `json:"args"` //查询、删除数据

	eg:[
			{"name":"productid","value":"PD000001"},
			{"name":"pageSize","value":"3"},
			{"name":"pageNumber","value":"2"}
		]

	Args:
		Name string `json:"name"`
		Value string `json:"value"`
}

*/

//从message对象的Args属性中获取到删除操作或者查询操作的条件，把这个条件返回给上一级接口
func GetArgsOfDelOrQuery(requestMsg *Message) ([]Args){
	var queryStringArr = requestMsg.Args
	return queryStringArr
}

//打印错误信息
func AddErrorToMessage(err string, tranCode string, returnMsg *Message)  {
	t := time.Now().Unix()
	fmt.Println(time.Unix(t, 0).String() + returnMsg.TraceNo + returnMsg.Orgid + returnMsg.Userid + ":" + err )
}

//从message对象中获取到data对象数组中的第一个元素
func GetOneDataFromMessage(requestMsg *Message) (string, error) {
	//取到data数组里面需要插入账本的产品JSON
	if len(requestMsg.Data) == 0{
		log.Error(ERRNMCC0012)
		var customErr = errors.New(ERRNMCC0012)
		return "", customErr
	}else {
		var productDataObj = requestMsg.Data[0] //data类型的对象

		var data Data
		data = productDataObj

		var productJSON = data.Content

		return productJSON, nil
	}

	return "", nil
}

//从message对象中获取到data数组中的所有对象元素
func GetDataFromMessage(requestMsg *Message) ([]Data, error) {
	//取到data数组里面需要插入账本的产品JSON
	if len(requestMsg.Data) == 0{
		log.Error(ERRNMCC0012)
		var customErr = errors.New(ERRNMCC0012)
		return nil, customErr
	}else {
		return requestMsg.Data, nil
	}

	return nil, nil
}

//通过couchDB的selector富查询来进行条件查询
func ConditionQuery(stub shim.ChaincodeStubInterface, queryString string) ([]string, error, string) {
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil{
		log.Error(ERRNMCC0006)
		return nil, err, ERRNMCC0006
	}
	defer resultsIterator.Close()

	var resArr []string

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil{
			log.Error(ERRNMCC0007)
			return nil, err, ERRNMCC0007
		}

		var buffer bytes.Buffer
		buffer.WriteString(string(queryResponse.Value))

		resArr = append(resArr, buffer.String())
	}

	return resArr, nil, ""
}

// 子串在字符串的字节位置
func UnicodeIndex(str,substr string) int {
	result := strings.Index(str,substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}

	return result
}

//截取字符串的方法
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//获取时间戳的方法 for example: 20171105164805 =>  2017-11-05 16:48:05
func GetTimeStamp() string {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	second := time.Now().Second()
	nsec := time.Now().Nanosecond()

	currentTime := time.Date(year, month, day, hour, minute, second, nsec, time.Local).String()

	rs := []rune(currentTime)
	res := string(rs[0: 10])

	canSplit := func (c rune) bool {return c=='-'}
	resArr := strings.FieldsFunc(res,canSplit)

	var result = ""
	for _, item := range resArr{
		result += item
	}

	resTime := string(rs[11: 19])

	canSplitTime := func (c rune) bool {return c==':'}
	resTimeArr := strings.FieldsFunc(resTime, canSplitTime)

	var resultTime = ""
	for _, item := range resTimeArr{
		resultTime += item
	}

	var resultTimeStamp = result+resultTime

	return resultTimeStamp
}

func ToBytes(entity interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// 通用的翻页查询
func PaginateQuery(stub shim.ChaincodeStubInterface, requestMsg *Message, returnMsg *Message) {
	// Fetch the query string and pagination information
	var args = GetArgsOfDelOrQuery(requestMsg)
	var pageSize int32
	var bookMark string
	var queryString string
	for _, argsItem := range args {
		switch argsItem.Name {
		case "pageSize":
			tmpPageSize, err := strconv.ParseInt(argsItem.Value, 10, 32)
			if err != nil {
				CombinationErrorObj(err.Error(), "PaginateQuery", ERRNMCC0009, returnMsg)
				ErrorMessage(returnMsg)
			}
			pageSize = int32(tmpPageSize)
			break
		case "bookMark":
			bookMark = argsItem.Value
			break
		case "queryString":
			queryString = argsItem.Value
			break
		}
	}
	if queryString == "" {
		tempQueryString, err := GetOneDataFromMessage(requestMsg)
		if err != nil {
			queryString = tempQueryString
		}
	}
	// 如果queryString为空
	if queryString == "" {
		queryString = "{\"selector\":{}}"
	}else{
		queryString = fmt.Sprintf("{\"selector\":%s}", queryString)
	}
	// 执行翻页查询
	log.Info(fmt.Sprintf("queryString: %s, pageSize: %d, bookMark: %s", queryString, pageSize, bookMark))
	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookMark)
	if err != nil {
		CombinationErrorObj(err.Error(), "PaginateQuery", "500", returnMsg)
		ErrorMessage(returnMsg)
	}else{
		// 定义返回的列表对象
		var rtList []interface{}
		defer resultsIterator.Close()
		// 循环处理数据
		for resultsIterator.HasNext() {
			// 获取一条记录信息
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				// 如果获取任意一条记录出错则返回错误信息
				CombinationErrorObj(err.Error(), "PaginateQuery", "ERR_NXT_REC", returnMsg)
				ErrorMessage(returnMsg)
			}
	
			var rtItem interface{}
			fmtErr := json.Unmarshal(queryResponse.Value, &rtItem)
			if fmtErr != nil {
				// 如果任意一条记录转换成JSON对象失败则直接返回
				CombinationErrorObj(err.Error(), "PaginateQuery", "ERR_FMT_REC", returnMsg)
				ErrorMessage(returnMsg)
			}
			rtList = append(rtList, rtItem)
		}
		var pagination McPagination
		pagination.CurBookMark = bookMark
		pagination.NxtBookMark = responseMetadata.Bookmark
		pagination.DataList = rtList
	
		returnMsg.RetCode = "000"
		returnMsg.RetObject = pagination
	}
}