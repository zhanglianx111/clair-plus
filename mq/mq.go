package mq

type Mqer interface {
	NewMq(name, tag, network, address string, db int) error
	NewConsumer(name string)
	//	GetMqClient() interface{}
	SendBytes(message []byte) bool
	SendString(message string) bool
}
