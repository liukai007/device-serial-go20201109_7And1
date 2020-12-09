package driver

import (
	"fmt"
	"strconv"
	"strings"
)

/*
strconv.ParseInt(s string, base int, bitSize int)
第二个参数base设置为0时，该函数会自动识别数字格式进行转换。
如果base为0，base的值会根据字符串s起始自动进行判断：如果字符串以0x开始，base=16；如果字符串以0开始base=8；以其他字符串开始base=10
*/
func Tran16StringTo10(content string, WholeOrseparate int) string {
	if WholeOrseparate == 0 {
		//使用空白字符进行分割
		strAarry := strings.Fields(content)
		bufTemp := make([]byte, len(strAarry))
		var i int
		for i = 0; i < len(strAarry); i++ {
			//16进制字符串转成10进制
			temp, _ := strconv.ParseInt("0x"+strAarry[i], 0, 16)
			fmt.Println("Tran16To10")
			fmt.Println(temp)
			bufTemp[i] = byte(temp)
		}
		str := string(bufTemp)
		return str
	} else {
		content = strings.Replace(content, " ", "", -1)
		//16进制转成10进制
		temp, _ := strconv.ParseInt("0x"+content, 0, 16)
		return string(temp)
	}

}

//十进制 转换成  16进制
func Tran10StringTo16(content string) string {
	//使用空白字符进行分割
	strAarry := strings.Fields(content)
	var resultString string
	var i int
	for i = 0; i < len(strAarry); i++ {
		//10进制字符串转成16进制
		temp, _ := strconv.ParseInt(strAarry[i], 0, 16)
		tmpStr := toHex(int(temp))
		if len(tmpStr) == 1 {
			tmpStr = "0" + tmpStr
		}
		if resultString == "" {
			resultString += tmpStr
		} else {
			resultString += " " + tmpStr
		}
	}
	return resultString
}

//十进制转换为16进制
func toHex(ten int) string {
	m := 0
	hex := make([]int, 0)
	for {
		m = ten % 16
		ten = ten / 16
		if ten == 0 {
			hex = append(hex, m)
			break
		}
		hex = append(hex, m)
	}
	var hexStr []string
	for i := len(hex) - 1; i >= 0; i-- {
		if hex[i] >= 10 {
			hexStr = append(hexStr, fmt.Sprintf("%c", 'A'+hex[i]-10))
		} else {
			hexStr = append(hexStr, fmt.Sprintf("%d", hex[i]))
		}
	}
	return strings.Join(hexStr, "")
}
