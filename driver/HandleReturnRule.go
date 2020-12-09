package driver

/*
处理返回规则
*/
type handleReturnRule struct {
	ReturnType               int               `json:"returnType"`               //0 表示字符串 1 int 2 float
	ReadTypeName             string            `json:"readTypeName"`             //读数类型名字
	RawResult                string            `json:"rawResult"`                //原始结果
	OtherReadTypeAndValue    map[string]string `json:"otherReadTypeAndValue"`    //
	ReadTypeReturnValueRange string            `json:"readTypeReturnValueRange"` ///返回值范围如果是当个的比如0,1 用逗号隔开，如果是1-100是范围
	Transcoding              int               `json:"transcoding"`              //是否需要转码 0 不需要转码  1表示10进制转换成16进制  2 表示16进制转换成10进制
	/*
	   *   EQ(0, "等于"),
	       NE(1, "不等于"),
	       LT(2, "小于"),
	       GT(3, "大于"),
	       LE(4, "小于等于"),
	       GE(5, "大于等于");
	       6 ==包含
	       7 ==不包含
	   * */
	JudgeSymbol int `json:"judgeSymbol"`

	//字符串 处理需要的字段，直接和返回值对比即可
	RightResult string `json:"rightResult"` //正确结果值  ==== 字符串 处理需要的字段，直接和返回值对比即可
	// 数字需要处理
	IsArrayType     int     `json:"isArrayType"`     //0 非数组类型  1数组类型
	GetSomeArray    string  `json:"getSomeArray"`    //比如 取第 4个，5个  ，4 5
	WholeOrseparate int     `json:"wholeOrseparate"` //whole=1 separate=0  1是整体  0是分开
	Scale           string  `json:"scale"`           //乘法因子
	AdditiveFactor  float64 `json:"additiveFactor"`  //加法因子
	//是否通过字符长度判断是否成功了
	IsLengthVerify int `json:"isLengthVerify"` //0 不通过长度验证 1 通过长度验证
}
