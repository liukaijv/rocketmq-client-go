package main

import (
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/liukaijv/rocketmq-client-go/core"
	"github.com/liukaijv/rocketmq-demo/api"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	namesrv = "localhost:9876"
	topic   = "topic2"
	groupId = "group1"
	amount  = 1

	rmq     = kingpin.New("rocketmq", "RocketMQ cmd tools")
	produce = rmq.Command("produce", "send messages to RocketMQ")
	consume = rmq.Command("consume", "consumes message from RocketMQ")
)

func main() {

	switch kingpin.MustParse(rmq.Parse(os.Args[1:])) {
	case produce.FullCommand():
		pConfig := &rocketmq.ProducerConfig{ClientConfig: rocketmq.ClientConfig{
			GroupID:    groupId,
			NameServer: namesrv,
			LogC: &rocketmq.LogConfig{
				Path:     "example",
				FileSize: 64 * 1 << 10,
				FileNum:  1,
				Level:    rocketmq.LogLevelDebug,
			},
		}}
		sendMessage(pConfig)
	case consume.FullCommand():
		cConfig := &rocketmq.PullConsumerConfig{ClientConfig: rocketmq.ClientConfig{
			GroupID:    groupId,
			NameServer: namesrv,
			LogC: &rocketmq.LogConfig{
				Path:     "example",
				FileSize: 64 * 1 << 10,
				FileNum:  1,
				Level:    rocketmq.LogLevelInfo,
			},
		}}

		consumeWithPull(cConfig, topic)
	}

}

func sendMessage(config *rocketmq.ProducerConfig) {
	producer, err := rocketmq.NewProducer(config)

	if err != nil {
		fmt.Println("create Producer failed, error:", err)
		return
	}

	err = producer.Start()
	if err != nil {
		fmt.Println("start producer error", err)
		return
	}
	defer producer.Shutdown()

	fmt.Printf("Producer: %s started... \n", producer)
	for i := 0; i < amount; i++ {

		registerRequest := api.RegisterServerRequest{
			Id:   proto.Int(1),
			Type: proto.String("edge"),
			Key:  proto.String("123456"),
			Name: proto.String("server1"),
			Port: proto.Int(0),
			Ip:   proto.String("127.0.0.1"),
		}
		data, err := proto.Marshal(&registerRequest)
		if err != nil {
			fmt.Println("proto.Marshal error:", err)
		}

		msg := &rocketmq.Message{Topic: topic, Body: data}

		result, err := producer.SendMessageSync(msg)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(hex.Dump(data))
		fmt.Printf("send message: %s result: %s\n", string(data), result)
	}
	fmt.Println("shutdown producer.")
}

func consumeWithPull(config *rocketmq.PullConsumerConfig, topic string) {

	consumer, err := rocketmq.NewPullConsumer(config)
	if err != nil {
		fmt.Printf("new pull consumer error:%s\n", err)
		return
	}

	err = consumer.Start()
	if err != nil {
		fmt.Printf("start consumer error:%s\n", err)
		return
	}
	defer consumer.Shutdown()

	mqs := consumer.FetchSubscriptionMessageQueues(topic)
	fmt.Printf("fetch subscription mqs:%+v\n", mqs)

	total, offsets, now := 0, map[int]int64{}, time.Now()

PULL:
	for {
		for _, mq := range mqs {
			pr := consumer.Pull(mq, "*", offsets[mq.ID], 32)
			total += len(pr.Messages)
			//fmt.Printf("pull %s, result:%+v\n", mq.String(), pr)

			for _, msg := range pr.Messages {

				fmt.Println(hex.Dump(msg.Body))

				var res = new(api.RegisterServerRequest)
				err = proto.Unmarshal(msg.Body, res)
				if err != nil {
					fmt.Println("proto.Unmarshal err", err)
					continue
				}

				fmt.Printf("pull %s, result:%+v\n", mq.String(), res)
				fmt.Println(res)

			}

			switch pr.Status {
			case rocketmq.PullNoNewMsg:
				break PULL
			case rocketmq.PullFound:
				fallthrough
			case rocketmq.PullNoMatchedMsg:
				fallthrough
			case rocketmq.PullOffsetIllegal:
				offsets[mq.ID] = pr.NextBeginOffset
			case rocketmq.PullBrokerTimeout:
				fmt.Println("broker timeout occur")
			}
		}
	}

	var timePerMessage time.Duration
	if total > 0 {
		timePerMessage = time.Since(now) / time.Duration(total)
	}
	fmt.Printf("total message:%d, per message time:%d\n", total, timePerMessage)
}
