package mq

import (
	"fmt"
	"testing"
	"time"
)

func TestNewMq(t *testing.T) {
	q := RedisMq{}
	err := q.NewMq("tasks", "my service", "tcp", "localhost:6379", 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(q)
	qq := q.GetMqClient()
	fmt.Println(qq)
	go q.NewConsumer("consumer")

	for i := 0; i < 10; i++ {
		q.Send(fmt.Sprintf("hello world %d", i))
	}
	time.Sleep(time.Second * time.Duration(1))
}
