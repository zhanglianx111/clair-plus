package mq

import (
	"encoding/json"
	"errors"
	"github.com/adjust/rmq"
	"github.com/astaxie/beego/logs"
	"github.com/zhanglianx111/clair-plus/clair"
	"github.com/zhanglianx111/clair-plus/models"
	"time"
)

type RedisMq struct {
	conn  rmq.Connection
	queue rmq.Queue
}

type Consumer struct {
	name  string
	count int
}

// new a message queue
func (r *RedisMq) NewMq(name, tag, network, address string, db int) error {
	c := rmq.OpenConnection(tag, network, address, db)
	if c == nil {
		logs.Error("open connection failed at ", network, "://", address)
		return errors.New("new a redis connection failed")
	}

	q := c.OpenQueue(name)
	if q == nil {
		logs.Error("open mq failed: ", name)
		c.Close()
		return errors.New("new a redis queue failed")
	}

	r.queue = q
	r.conn = c
	return nil
}

// get mq client
func (r RedisMq) GetMqClient() rmq.Queue {
	return r.queue
}

// add a consumer into mq
func (r RedisMq) NewConsumer(name string) {
	r.queue.StartConsuming(1000, 500*time.Millisecond)
	consumer := Consumer{
		name:  name,
		count: 0,
	}
	r.queue.AddConsumer(name, &consumer)
	select {}
}

// do a message task
func (consumer *Consumer) Consume(message rmq.Delivery) {
	var image models.Image

	if err := json.Unmarshal([]byte(message.Payload()), &image); err != nil {
		message.Reject()
		return
	}

	consumer.count++

	logs.Info("get message: ", image, "count: ", consumer.count)

	beginTime := time.Now()
	scanedLayer, err := clair.GetClairHandler().ScanAndGetFeatures(image.Repo, image.Tag)
	if err != nil {
		logs.Error("扫描images失败:", err)
		return
	} else {
		message.Ack()
	}
	logs.Info(scanedLayer)
	// send vnlnerabilites to somewhere
	elapsed := time.Since(beginTime)
	logs.Info("执行时间:", elapsed)
	/*
		s.Data["json"] = scanedLayer
		s.Data["json"] = result
		s.ServeJSON()
	*/
}

// add a string message into mq
func (r *RedisMq) SendString(message string) bool {
	return r.queue.Publish(message)
}

// add a byte message into mq
func (r *RedisMq) SendBytes(message []byte) bool {
	return r.queue.PublishBytes(message)
}
