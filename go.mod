module github.com/edgexfoundry/device-serial-go

go 1.12

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0 // indirect
	github.com/edgexfoundry/device-sdk-go v1.1.1
	github.com/edgexfoundry/go-mod-core-contracts v0.1.31
	github.com/shopspring/decimal v1.2.0
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20181029021203-45a5f77698d3
	golang.org/x/net => github.com/golang/net v0.0.0-20190228165749-92fc7df08ae7
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180823144017-11551d06cbcc
)
