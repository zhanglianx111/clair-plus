package mq

import (
	"encoding/json"
	"errors"
	"github.com/adjust/rmq"
	"github.com/astaxie/beego/logs"
	"github.com/zhanglianx111/clair-plus/clair"
	"github.com/zhanglianx111/clair-plus/models"
	"time"
	"github.com/coreos/clair/api/v1"
	"github.com/astaxie/beego/httplib"
	"strings"
	"github.com/astaxie/beego"
)

type RedisMq struct {
	conn  rmq.Connection
	queue rmq.Queue
}

type Consumer struct {
	name  string // consumer's name
	count int    // count of message consumed
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
	r.queue.StartConsuming(1000, 100*time.Millisecond)
	consumer := Consumer{
		name:  name,
		count: 0,
	}

	defer r.queue.Close()
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
	logs.Debug(scanedLayer)

	// TODO
	elapsed := time.Since(beginTime)
	logs.Info("执行时间:", elapsed)

	// send vnlnerabilites to somewhere
	//现在发从给测试程序
	sendStr := sendStruct{
		layer: scanedLayer.Layer,
		error: scanedLayer.Error,
		usedTime: elapsed.String(),
	}
	sendResult(sendStr, image)
}

// add a string message into mq
func (r *RedisMq) SendString(message string) bool {
	return r.queue.Publish(message)
}

// add a byte message into mq
func (r *RedisMq) SendBytes(message []byte) bool {
	return r.queue.PublishBytes(message)
}

func sendResult(sendStr sendStruct, image models.Image) {

	webUrl := beego.AppConfig.String("webURL")

	spl := strings.Split(image.Repo, "/")
	namespace := spl[0]
	imageName := spl[1]

	sendURL :=  webUrl + "/v1/clair/" + "registry/hub.hcpaas.com/namespace/" + namespace + "/image/" + imageName + "/tag/" + image.Tag + "/imageReport"

	req := httplib.Put(sendURL)

	req, err := req.JSONBody(sendStr)
	if err != nil {
		logs.Error("转换失败:", err)
	}
	
	logs.Warning(sendStr)

	req.Header("Content-Type", "application/json;charset=utf-8")

	resp, err := req.DoRequest()
	if err != nil {
		logs.Error("向web port发送put请求失败:", err)
	}
	if resp.StatusCode != 200 {
		logs.Error("向web port发送put请求失败:", resp.Status)
	}
	logs.Debug("向web port发送成功")
}

type sendStruct struct {
	layer *v1.Layer `json:"layer"`
	error *v1.Error `json:"error"`
	usedTime string `json:"usedTime"`
}