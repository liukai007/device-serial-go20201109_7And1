package driver

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func handleFloatValue(byte1 []byte, getValueArray string, Scale string, AdditiveFactor float64, WholeOrseparate int) float64 {
	temp1 := strings.Fields(getValueArray)
	var i int
	var valueString string
	for i = 0; i < len(temp1); i++ {
		i1, _ := strconv.ParseInt(temp1[i], 0, 8)
		h := fmt.Sprintf("%X", byte1[i1])
		if len(h) == 1 {
			h = "0" + h
		}
		valueString = valueString + h
	}
	fmt.Println("要处理的值：" + valueString)
	valueInt, _ := strconv.ParseInt(valueString, 16, 16)
	fmt.Printf("integer%v", valueInt)
	var test decimal.Decimal
	test = decimal.NewFromFloat(float64(valueInt))
	//乘法因子
	if Scale != "" {
		scalefloat64, _ := strconv.ParseFloat(Scale, 64)
		test = test.Mul(decimal.NewFromFloat(scalefloat64))
	}
	if AdditiveFactor != 0.0 {
		test = test.Add(decimal.NewFromFloat(AdditiveFactor))
	}
	value, _ := test.Float64()
	return value
}

type Protocols struct {
	SerialPort string `json:"serialPort"`
	BaudRate   string `json:"baudRate"`
	DataBits   string `json:"dataBit"`
	StopBits   string `json:"stopBit"`
	Parity     string `json:"parityBit"`
}

//通过http得到该设备的链接属性
func getProperty(ipAddress string, ipPort string, equipmentName string) Protocols {
	url := "http://" + ipAddress + ":" + ipPort + "/portType/equipmentEnName/" + equipmentName
	resp, _ := http.Get(url)
	protocols := Protocols{}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	//字符串替换
	content := string(body)
	content = strings.Replace(content, "{\"data\":\"", "", -1)
	content = strings.Replace(content, "\"}\"}", "\"}", -1)
	content = strings.Replace(content, "\\\"", "\"", -1)
	fmt.Println(content)
	json.Unmarshal([]byte(content), &protocols)
	return protocols
}

//map[string]string 两个合并成一个map
func mapPutAll(map1 map[string]string, map2 map[string]string) map[string]string {
	mapText1, _ := json.Marshal(map1)
	mapText2, _ := json.Marshal(map2)
	mapText1Str := strings.Replace(string(mapText1), "}", ",", -1)
	mapText2Str := strings.Replace(string(mapText2), "{", "", -1)
	mapText1Str = mapText1Str + mapText2Str
	mapText1Str = strings.Replace(mapText1Str, "\"\":\"\",", "", -1)
	var map3 map[string]string
	json.Unmarshal([]byte(mapText1Str), &map3)
	return map3
}

//! 字符串转数字
func HF_Atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

//! 数字转字符串
func HF_Itoa(s int) string {
	num := strconv.Itoa(s)
	return num
}

//! 字符串转float32
func HF_Atof(s string) float32 {
	num, _ := strconv.ParseFloat(s, 32)
	return float32(num)
}

//! 字符串转float64
func HF_Atof64(s string) float64 {
	num, _ := strconv.ParseFloat(s, 64)
	return num
}

//!float转换为字符串
/*
HF_Atos(value,64)
*/
func HF_Atos(s float64, f int) string {
	num := strconv.FormatFloat(s, 'E', -1, f)
	return num
}
