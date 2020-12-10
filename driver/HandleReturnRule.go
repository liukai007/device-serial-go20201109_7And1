package driver

/*
处理返回规则
*/
type handleReturnRule struct {
	ReturnType                int               `json:"returnType"`                //0 表示字符串 1 int 2 float
	ReadTypeName              string            `json:"readTypeName"`              //读数类型名字
	RawResult                 string            `json:"rawResult"`                 //原始结果
	OtherReadTypeAndValue     map[string]string `json:"otherReadTypeAndValue"`     //其他属性值，可以手动填写进去
	ReadTypeReturnValueRange  string            `json:"readTypeReturnValueRange"`  ///返回值范围如果是当个的比如0,1 用逗号隔开，如果是1-100是范围
	Transcoding               int               `json:"transcoding"`               //是否需要转码 0 不需要转码  1表示10进制转换成16进制  2 表示16进制转换成10进制
	RegularExpression         string            `json:"regularExpression"`         //正则表达式 主要针对结果进行处理的
	RegularExpressionWhichOne int               `json:"regularExpressionWhichOne"` //正则表达式 如果好几个，默认取第一个，但是可以填写 0表示第一个，1表示第二个，2表示第三个
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
	JudgeSymbol int `json:"judgeSymbol"` //进行条件判断，比如字符串如果相同就是

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
