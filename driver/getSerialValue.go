package driver

import (
	"fmt"
	"github.com/tarm/serial"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*
 串口执行的时候，有时候执行一次， 因为前面有一些遗留的返回数据，导致返回异常，所以要多执行1次才ok，特定设备的偶尔性故障解决方案
*/
func getValuebyCmdStringTimes(CmdContent string, serialPort string, baudRate string, dataBits string, stopBits string, parity string, times int) (string, []byte, error) {
	if times > 1 {
		for j := 0; j < times-1; j++ {
			getValuebyCmdString(CmdContent, serialPort, baudRate, dataBits, stopBits, parity)
		}
	}
	return getValuebyCmdString(CmdContent, serialPort, baudRate, dataBits, stopBits, parity)
}

/*
串口连接，并且执行获取值
*/
func getValuebyCmdString(CmdContent string, serialPort string, baudRate string, dataBits string, stopBits string, parity string) (string, []byte, error) {
	fmt.Println("进入getValuebyCmdString")
	baudInt, err := strconv.Atoi(baudRate)
	if err != nil {
		return "Baud错误", nil, err
	}

	config := &serial.Config{
		Name:        serialPort,
		Baud:        baudInt,
		ReadTimeout: 3 * time.Second,
	}
	fmt.Println("打开串口" + serialPort)
	s, err := serial.OpenPort(config)
	defer s.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Println("串口被占用，沉睡5秒中")
		time.Sleep(5 * time.Second)
		exec.Command("fuser -k /dev/ttyS2")
		s, err = serial.OpenPort(config)
		if err != nil {
			return "", nil, err
		}
	}
	fmt.Println("连接成功" + serialPort)
	//字符串的十六进制转成十进制
	serialStrAarry := strings.Fields(CmdContent)
	long1 := len(serialStrAarry)
	bufTemp := make([]byte, long1)
	var i int
	for i = 0; i < long1; i++ {
		temp, _ := strconv.ParseInt("0x"+serialStrAarry[i], 0, 16)
		bufTemp[i] = byte(temp)
	}
	s.Write(bufTemp)
	buf := make([]byte, 100)
	n, err := s.Read(buf)
	fmt.Println(buf[:n])
	//把buf[:n]转成字符串
	str := convertByteToString(buf[:n])
	fmt.Println("转化完成的字符串：" + str)
	return str, buf[:n], err
}
