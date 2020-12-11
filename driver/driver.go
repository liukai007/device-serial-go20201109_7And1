// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example implementation of
// ProtocolDriver interface.
//
package driver

import (
	"encoding/json"
	"fmt"
	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SimpleDriver struct {
	lc      logger.LoggingClient
	asyncCh chan<- *dsModels.AsyncValues
}

//初始化这个设备的执行协议
// service.
func (s *SimpleDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	s.lc = lc
	s.asyncCh = asyncCh
	fmt.Println("Initialize success!!!!")
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *SimpleDriver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	res = make([]*dsModels.CommandValue, len(reqs))
	now := time.Now().UnixNano()
	fmt.Println("DeviceName: " + deviceName)
	//步骤1 得到串口的相关信息
	//url := "http://192.168.2.10:8989/portType/equipmentEnName/7he1"
	var ip string     //ip地址
	var ipPort string //ip端口
	for k, v := range protocols {
		fmt.Println("key he value")
		fmt.Println(k, v)
		ip = v["ip"]
		ipPort = v["ipPort"]
	}
	protocol1 := getProperty(ip, ipPort, deviceName)
	fmt.Println(protocol1)
	//步骤2 得到返回的结果
	//得到命令字符串
	attributes1 := attributes{}
	tmpMap := reqs[0].Attributes
	if tmpMap != nil {
		attributes1.cmdContent = tmpMap["cmdContent"]
		attributes1.transcoding = tmpMap["transcoding"]
		attributes1.handleReturnRules = tmpMap["handleReturnRules"]
		attributes1.executionTimes = HF_Atoi(tmpMap["executionTimes"])
		attributes1.regularExpression = tmpMap["regularExpression"]
		attributes1.regularExpressionWhichOne = HF_Atoi(tmpMap["regularExpressionWhichOne"])
	}
	var handleReturnRuleList []handleReturnRule

	if attributes1.handleReturnRules != "" {
		handleReturnRules := attributes1.handleReturnRules
		json.Unmarshal([]byte(handleReturnRules), &handleReturnRuleList)
	}
	//多执行一次保证获取的字段正确
	//参数 命令， 串口号，波特率，数据位，停止位，奇偶校验位
	returnResult, by1, err := getValuebyCmdStringTimes(attributes1.cmdContent, protocol1.SerialPort, protocol1.BaudRate, protocol1.DataBits, protocol1.StopBits, protocol1.Parity, attributes1.executionTimes)
	/*************没有连接成功****开始*******/
	if err != nil {
		var cv1 *dsModels.CommandValue
		for i, req := range reqs {
			var returnStringMap map[string]string /*创建集合 */
			returnStringMap = make(map[string]string)
			returnStringMap["equipmentNameEn"] = deviceName
			returnStringMap["runningStatus"] = HF_Itoa(2)
			valueString, _ := json.Marshal(returnStringMap)
			cv1 = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
			res[i] = cv1
		}
		return res, nil
	}
	/*************没有连接成功****结束*****  **/
	fmt.Println("返回结果：" + returnResult)
	/************************步骤3 - 开始 ************************************/
	//是否需要转码
	//处理得到返回值
	if attributes1.transcoding == "1" {
		returnResult = Tran10StringTo16(returnResult)
		fmt.Println("转码返回结果：" + returnResult)
	} else if attributes1.transcoding == "2" {
		returnResult = Tran16StringTo10(returnResult, 0)
	}
	//是否使用正则表达式
	if attributes1.regularExpression != "" {
		fmt.Println("正则前。。。")
		fmt.Println(attributes1.regularExpression)
		fmt.Println(returnResult)
		if attributes1.regularExpressionWhichOne == 0 {
			reg := regexp.MustCompile(attributes1.regularExpression)
			returnResult = reg.FindString(returnResult)
		} else {
			reg := regexp.MustCompile(attributes1.regularExpression)
			MEIds := reg.FindAllString(returnResult, 3)
			for i := 0; i < len(MEIds); i++ {
				if i == attributes1.regularExpressionWhichOne {
					returnResult = MEIds[i]
					break
				}
			}
		}
		fmt.Println("正则后")
	}

	//res = make([]*dsModels.CommandValue, len(reqs))
	if len(handleReturnRuleList) == 1 {
		fmt.Println("follow1")
		for j := 0; j < len(handleReturnRuleList); j++ {
			var returnStringMap map[string]string /*创建集合 */
			returnStringMap = make(map[string]string)
			returnStringMap["equipmentNameEn"] = deviceName
			/*如果返回数组是字符型*/
			if handleReturnRuleList[j].ReturnType == 0 {
				/*
					EQ(0, "等于"),
					NE(1, "不等于"),
					LT(2, "小于"),
					GT(3, "大于"),
					LE(4, "小于等于"),
					GE(5, "大于等于");
					6 ==包含
					7 ==不包含
				*/
				if handleReturnRuleList[j].JudgeSymbol == 0 {
					if returnResult == handleReturnRuleList[j].RightResult {
						returnStringMap = mapSuccess1(returnStringMap)
						//处理字符串根据 逗号  转换成数组
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							fmt.Println(handleReturnRuleList[j].ReadTypeName)
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
					//不等于
				} else if handleReturnRuleList[j].JudgeSymbol == 1 {
					if returnResult != handleReturnRuleList[j].RightResult {
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
						if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
							returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
					//包含
				} else if handleReturnRuleList[j].JudgeSymbol == 6 {
					if strings.Contains(returnResult, handleReturnRuleList[j].RightResult) {
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
						if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
							returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
					//不包含
				} else if handleReturnRuleList[j].JudgeSymbol == 7 {
					if !strings.Contains(returnResult, handleReturnRuleList[j].RightResult) {
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
						if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
							returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
				}

				var cv *dsModels.CommandValue
				for i, req := range reqs {
					if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
						returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
					}
					valueString, _ := json.Marshal(returnStringMap)
					cv = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
					res[i] = cv
				}

				//一些数值型的计算----float
			} else if handleReturnRuleList[j].ReturnType == 2 {
				returnStringMap = mapSuccess1(returnStringMap)
				value := handleFloatValue(by1, handleReturnRuleList[j].GetSomeArray, handleReturnRuleList[j].Scale, handleReturnRuleList[j].AdditiveFactor, handleReturnRuleList[j].WholeOrseparate)
				fmt.Printf("%f", value)
				if handleReturnRuleList[j].ReadTypeName != "" && handleReturnRuleList[j].RightResult != "" {
					/*
						func judgeSymbolAndRightResult(readTypeName string, value int, JudgeSymbol int, RightResult string) map[string]string {
					*/
					result := judgeSymbolAndRightResultFloat(handleReturnRuleList[j].ReadTypeName, value, handleReturnRuleList[j].JudgeSymbol, handleReturnRuleList[j].RightResult)
					var cv1 *dsModels.CommandValue
					for i, req := range reqs {
						var returnStringMap map[string]string /*创建集合 */
						returnStringMap = make(map[string]string)
						returnStringMap["equipmentNameEn"] = deviceName
						returnStringMap = mapPutAll(returnStringMap, result)
						valueString, _ := json.Marshal(returnStringMap)
						cv1 = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
						res[i] = cv1
					}
					return res, nil
				}

				/*****************判断范围--开始************************/
				if handleReturnRuleList[j].ReadTypeReturnValueRange != "" {
					string1 := handleReturnRuleList[j].ReadTypeReturnValueRange
					if strings.Contains(string1, "-") {
						list := strings.Split(string1, "-")
						if len(list) == 2 {
							low := HF_Atof64(list[0])
							high := HF_Atof64(list[1])
							if value > high || value <= low {
								var cv1 *dsModels.CommandValue
								for i, req := range reqs {
									var returnStringMap map[string]string /*创建集合 */
									returnStringMap = make(map[string]string)
									returnStringMap["equipmentNameEn"] = deviceName
									returnStringMap["runningStatus"] = HF_Itoa(2)
									returnStringMap = mapSuccess0(returnStringMap)
									valueString, _ := json.Marshal(returnStringMap)
									cv1 = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
									res[i] = cv1
								}
								return res, nil
							}
						}
					}
				}
				/*****************判断范围--结束************************/
				var cv *dsModels.CommandValue
				for i, req := range reqs {
					returnStringMap[handleReturnRuleList[j].ReadTypeName] = strconv.FormatFloat(float64(value), 'f', 6, 64)
					if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
						returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
					}
					valueString, _ := json.Marshal(returnStringMap)
					cv = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
					res[i] = cv
				}

				//一些数值型的计算----int
			} else if handleReturnRuleList[j].ReturnType == 1 {
				returnStringMap = mapSuccess1(returnStringMap)
				value := handleFloatValue(by1, handleReturnRuleList[j].GetSomeArray, handleReturnRuleList[j].Scale, handleReturnRuleList[j].AdditiveFactor, handleReturnRuleList[j].WholeOrseparate)
				fmt.Printf("%f", value)
				if handleReturnRuleList[j].ReadTypeName != "" && handleReturnRuleList[j].RightResult != "" {
					result := judgeSymbolAndRightResultInt(handleReturnRuleList[j].ReadTypeName, int(value), handleReturnRuleList[j].JudgeSymbol, handleReturnRuleList[j].RightResult)
					var cv1 *dsModels.CommandValue
					for i, req := range reqs {
						var returnStringMap map[string]string /*创建集合 */
						returnStringMap = make(map[string]string)
						returnStringMap["equipmentNameEn"] = deviceName
						returnStringMap = mapPutAll(returnStringMap, result)
						valueString, _ := json.Marshal(returnStringMap)
						cv1 = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
						res[i] = cv1
					}
					return res, nil
				}

				/*****************判断范围--开始************************/
				if handleReturnRuleList[j].ReadTypeReturnValueRange != "" {
					string1 := handleReturnRuleList[j].ReadTypeReturnValueRange
					if strings.Contains(string1, "-") {
						list := strings.Split(string1, "-")
						if len(list) == 2 {
							low := HF_Atof64(list[0])
							high := HF_Atof64(list[1])
							if value > high || value <= low {
								var cv1 *dsModels.CommandValue
								for i, req := range reqs {
									var returnStringMap map[string]string /*创建集合 */
									returnStringMap = make(map[string]string)
									returnStringMap["equipmentNameEn"] = deviceName
									returnStringMap["runningStatus"] = HF_Itoa(2)
									returnStringMap = mapSuccess0(returnStringMap)
									valueString, _ := json.Marshal(returnStringMap)
									cv1 = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
									res[i] = cv1
								}
								return res, nil
							}
						}
					}
				}
				/*****************判断范围--结束************************/
				var cv *dsModels.CommandValue
				for i, req := range reqs {
					returnStringMap[handleReturnRuleList[j].ReadTypeName] = HF_Atos(value, 64)
					if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
						returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
					}
					valueString, _ := json.Marshal(returnStringMap)
					cv = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
					res[i] = cv
				}
			}

		}
	}
	//如果 多个取数规则的话，走下面的这一路
	if len(handleReturnRuleList) > 1 {
		fmt.Println("多个规则处理")
		var returnStringMap map[string]string /*创建集合 */
		returnStringMap = make(map[string]string)
		returnStringMap["equipmentNameEn"] = deviceName
		for j := 0; j < len(handleReturnRuleList); j++ {
			/*如果返回数组是字符型*/
			if handleReturnRuleList[j].ReturnType == 0 {
				/*
					EQ(0, "等于"),
					NE(1, "不等于"),
					LT(2, "小于"),
					GT(3, "大于"),
					LE(4, "小于等于"),
					GE(5, "大于等于");
					6 ==包含
					7 ==不包含
				*/
				if handleReturnRuleList[j].JudgeSymbol == 0 {
					if returnResult == handleReturnRuleList[j].RightResult {
						returnStringMap = mapSuccess1(returnStringMap)
						//处理字符串根据 逗号  转换成数组
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							fmt.Println(handleReturnRuleList[j].ReadTypeName)
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
					//不等于
				} else if handleReturnRuleList[j].JudgeSymbol == 1 {
					if returnResult != handleReturnRuleList[j].RightResult {
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
						if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
							returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
					//包含
				} else if handleReturnRuleList[j].JudgeSymbol == 6 {
					if strings.Contains(returnResult, handleReturnRuleList[j].RightResult) {
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
						if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
							returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
					//不包含
				} else if handleReturnRuleList[j].JudgeSymbol == 7 {
					if !strings.Contains(returnResult, handleReturnRuleList[j].RightResult) {
						rtSz := strings.Split(handleReturnRuleList[j].ReadTypeReturnValueRange, ",")
						if len(rtSz) == 1 {
							returnStringMap[handleReturnRuleList[j].ReadTypeName] = rtSz[0]
						}
						if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
							returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
						}
					} else {
						returnStringMap = mapSuccess0(returnStringMap)
					}
				}
				if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
					returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
				}

				//一些数值型的计算----float
			} else if handleReturnRuleList[j].ReturnType == 2 {
				returnStringMap = mapSuccess1(returnStringMap)
				value := handleFloatValue(by1, handleReturnRuleList[j].GetSomeArray, handleReturnRuleList[j].Scale, handleReturnRuleList[j].AdditiveFactor, handleReturnRuleList[j].WholeOrseparate)
				fmt.Printf("%f", value)
				if handleReturnRuleList[j].ReadTypeName != "" && handleReturnRuleList[j].RightResult != "" {
					/*
						func judgeSymbolAndRightResult(readTypeName string, value int, JudgeSymbol int, RightResult string) map[string]string {
					*/
					result := judgeSymbolAndRightResultFloat(handleReturnRuleList[j].ReadTypeName, value, handleReturnRuleList[j].JudgeSymbol, handleReturnRuleList[j].RightResult)
					returnStringMap = mapPutAll(returnStringMap, result)
				}

				/*****************判断范围--开始************************/
				if handleReturnRuleList[j].ReadTypeReturnValueRange != "" {
					string1 := handleReturnRuleList[j].ReadTypeReturnValueRange
					if strings.Contains(string1, "-") {
						list := strings.Split(string1, "-")
						if len(list) == 2 {
							low := HF_Atof64(list[0])
							high := HF_Atof64(list[1])
							if value > high || value <= low {
								var cv1 *dsModels.CommandValue
								for i, req := range reqs {
									var returnStringMap map[string]string /*创建集合 */
									returnStringMap = make(map[string]string)
									returnStringMap["equipmentNameEn"] = deviceName
									returnStringMap["runningStatus"] = HF_Itoa(2)
									returnStringMap = mapSuccess0(returnStringMap)
									valueString, _ := json.Marshal(returnStringMap)
									cv1 = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
									res[i] = cv1
								}
								return res, nil
							}
						}
					}
				}
				/*****************判断范围--结束************************/
				returnStringMap[handleReturnRuleList[j].ReadTypeName] = strconv.FormatFloat(float64(value), 'f', 6, 64)
				if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
					returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
				}
				//一些数值型的计算----int
			} else if handleReturnRuleList[j].ReturnType == 1 {
				returnStringMap = mapSuccess1(returnStringMap)
				value := handleFloatValue(by1, handleReturnRuleList[j].GetSomeArray, handleReturnRuleList[j].Scale, handleReturnRuleList[j].AdditiveFactor, handleReturnRuleList[j].WholeOrseparate)
				fmt.Printf("%f", value)
				if handleReturnRuleList[j].ReadTypeName != "" && handleReturnRuleList[j].RightResult != "" {
					result := judgeSymbolAndRightResultInt(handleReturnRuleList[j].ReadTypeName, int(value), handleReturnRuleList[j].JudgeSymbol, handleReturnRuleList[j].RightResult)
					returnStringMap = mapPutAll(returnStringMap, result)
				}

				/*****************判断范围--开始************************/
				if handleReturnRuleList[j].ReadTypeReturnValueRange != "" {
					string1 := handleReturnRuleList[j].ReadTypeReturnValueRange
					if strings.Contains(string1, "-") {
						list := strings.Split(string1, "-")
						if len(list) == 2 {
							low := HF_Atof64(list[0])
							high := HF_Atof64(list[1])
							if value > high || value <= low {
								var cv1 *dsModels.CommandValue
								for i, req := range reqs {
									var returnStringMap map[string]string /*创建集合 */
									returnStringMap = make(map[string]string)
									returnStringMap["equipmentNameEn"] = deviceName
									returnStringMap["runningStatus"] = HF_Itoa(2)
									returnStringMap = mapSuccess0(returnStringMap)
									valueString, _ := json.Marshal(returnStringMap)
									cv1 = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
									res[i] = cv1
								}
								return res, nil
							}
						}
					}
				}
				/*****************判断范围--结束************************/
				returnStringMap[handleReturnRuleList[j].ReadTypeName] = HF_Atos(value, 64)
				if len(handleReturnRuleList[j].OtherReadTypeAndValue) > 0 {
					returnStringMap = mapPutAll(returnStringMap, handleReturnRuleList[j].OtherReadTypeAndValue)
				}
			}

		}
		//最后返回的东西
		var cv *dsModels.CommandValue
		for i, req := range reqs {
			valueString, _ := json.Marshal(returnStringMap)
			cv = dsModels.NewStringValue(req.DeviceResourceName, now, string(valueString))
			res[i] = cv
		}
		return res, nil
	}

	return res, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *SimpleDriver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {
	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *SimpleDriver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if s.lc != nil {
		s.lc.Debug(fmt.Sprintf("SimpleDriver.Stop called: force=%v", force))
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *SimpleDriver) AddDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *SimpleDriver) UpdateDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *SimpleDriver) RemoveDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	s.lc.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}

func judgeSymbolAndRightResultInt(readTypeName string, value int, JudgeSymbol int, RightResult string) map[string]string {
	var result map[string]string
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
	if JudgeSymbol == 0 {
		if value == HF_Atoi(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*NE(1, "不等于"),*/
	if JudgeSymbol == 1 {
		if value != HF_Atoi(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*LT(2, "小于"),*/
	if JudgeSymbol == 2 {
		if value < HF_Atoi(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*GT(3, "大于"),*/
	if JudgeSymbol == 3 {
		if value > HF_Atoi(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*LE(4, "小于等于"),*/
	if JudgeSymbol == 4 {
		if value <= HF_Atoi(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*GE(5, "大于等于")*/
	if JudgeSymbol == 5 {
		if value >= HF_Atoi(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	return mapSuccess0(result)
}

func judgeSymbolAndRightResultFloat(readTypeName string, value float64, JudgeSymbol int, RightResult string) map[string]string {
	var result map[string]string
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
	if JudgeSymbol == 0 {
		if value == HF_Atof64(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*NE(1, "不等于"),*/
	if JudgeSymbol == 1 {
		if value != HF_Atof64(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*LT(2, "小于"),*/
	if JudgeSymbol == 2 {
		if value < HF_Atof64(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*GT(3, "大于"),*/
	if JudgeSymbol == 3 {
		if value > HF_Atof64(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*LE(4, "小于等于"),*/
	if JudgeSymbol == 4 {
		if value <= HF_Atof64(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}
	/*GE(5, "大于等于")*/
	if JudgeSymbol == 5 {
		if value >= HF_Atof64(RightResult) {
			return resultReadTypeName(readTypeName, "1")
		} else {
			return resultReadTypeName(readTypeName, "0")
		}
	}

	return mapSuccess0(result)
}
func resultReadTypeName(readTypeName string, readTypeValue string) map[string]string {
	var result map[string]string
	result[readTypeName] = readTypeValue
	return mapSuccess1(result)
}
