// -------------------------------------------
// @file      : producer.go
// @author    : bo cai
// @contact   : caibo923@gmail.com
// @time      : 2025/1/6 下午5:51
// -------------------------------------------

package main

import (
	"github.com/IBM/sarama"
	"github.com/caibo86/cberrors"
	"github.com/caibo86/logger"
	"time"
)

func main() {
	logger.Init(
		logger.SetIsOpenFile(false),
		logger.SetIsRedirectErr(false))
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		cberrors.Panic("create producer failed, err:%s", err)
	}
	defer func() {
		_ = producer.Close()
	}()
	// 发布消息
	msg := &sarama.ProducerMessage{
		Topic: "art",
		Value: sarama.StringEncoder("test message"),
	}
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				logger.Errorf("send message failed, err:%s", err)
			} else {
				logger.Infof("send message success, partition:%d, offset:%d", partition, offset)
			}
		}
	}

}
