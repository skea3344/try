// -------------------------------------------
// @file      : consumer.go
// @author    : bo cai
// @contact   : caibo923@gmail.com
// @time      : 2025/1/7 上午10:43
// -------------------------------------------

package main

import (
	"github.com/IBM/sarama"
	"github.com/caibo86/cberrors"
	"github.com/caibo86/logger"
)

func main() {
	logger.Init(
		logger.SetIsOpenFile(false),
		logger.SetIsRedirectErr(false),
	)
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"172.19.0.2:9092"}, config)
	if err != nil {
		cberrors.PanicWrap(err)
	}
	defer func() {
		_ = consumer.Close()
	}()
	partitionConsumer, err := consumer.ConsumePartition("art", 0, sarama.OffsetNewest)
	if err != nil {
		logger.Errorf("consume partition failed, err:%s", err)
		return
	}
	defer func() {
		_ = partitionConsumer.Close()
	}()
	for msg := range partitionConsumer.Messages() {
		logger.Infof("receive message, topic:%s, partition:%d, offset:%d, value:%s", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
	}
}
