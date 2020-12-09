package driver

type attributes struct {
	cmdContent        string //命令字符串
	handleReturnRules string //是一个对象处理返回规则
	transcoding       string //是否需要转码 0 不需要转码  1表示10进制转换成16进制
}
