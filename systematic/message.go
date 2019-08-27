package systematic

type Message struct {
	Id  string `json:"id"`
	Channel  string `json:"channel"`
	TranCode  string `json:"tranCode"`
	Type  string `json:"type"`
	DataType  string `json:"dataType"`
	TranDate  string `json:"tranDate"`
	TraceNo  string `json:"traceNo"` //
	RetCode  string `json:"retCode"` //处理过程中如果发生错误，把该值置为:111   处理过程成功，把该值置为:000
	Debug  string `json:"debug"`
	Userid  string `json:"userid"`
	Orgid  string `json:"orgid"`
	HasGeneric bool `json:"hasGeneric"`
	RetObject interface{} `json:"retObject"`

	Data  []Data `json:"data"`	     //新增、修改数据
	Errors  []Error `json:"errors"`
	Args   []Args `json:"args"` //查询、删除数据    eg:[{"name":"productid","value":"PD000001"}]
}

type Args struct { //查询条件
	Name string `json:"name"`
	Value string `json:"value"`
}

type Data struct {
	Content  string `json:"content"`
	DataType  string `json:"dataType"`
	FormatType string `json:"formatType"`
}

type Error struct {
	Statement  string `json:"statement"`
	Function  string `json:"function"`
	Code  string `json:"code"`
}

// 初始化返回报文对象
func InitMessage(id string, channel string, trancode string, trantime string, requestMsg *Message) (msg Message) {
	var returnMsg Message
	returnMsg.Id = id
	returnMsg.Channel = channel
	returnMsg.TranCode = trancode
	returnMsg.Orgid = requestMsg.Orgid
	returnMsg.Userid = requestMsg.Userid
	returnMsg.TranDate = ""
	returnMsg.Type = "1" //0是请求报文，1是返回报文
	returnMsg.RetCode = "000"

	return returnMsg
}

func ErrorMessage(returnMsg *Message) {
	returnMsg.RetCode = "111"
}

//分页使用的结构体
type McPagination struct {
	CurBookMark string        `json:"curBookMark"`
	NxtBookMark string        `json:"nxtBookMark"` //保存页码hash值
	PreBookMark string        `json:"preBookMark"`
	PageSize    string        `json:"pageSize"`
	DataList    interface{}   `json:"dataList"` //保存当前页的数据内容
}