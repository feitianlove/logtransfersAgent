package main

import (
	golibKafak "github.com/feitianlove/golib/kafka"
	"github.com/feitianlove/logtransfersAgent/config"
	"github.com/feitianlove/logtransfersAgent/kafka"
	"github.com/feitianlove/logtransfersAgent/logger"
	"github.com/feitianlove/logtransfersAgent/tailf"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// 初始化config
	cfg, err := config.NewConfig("../etc/logtransferAgent.conf")
	if err != nil {
		panic(err)
	}
	// 初始化 logger
	err = logger.InitLog(cfg)
	if err != nil {
		panic(err)
	}
	//初始化kafka
	kafkaProduct, err := kafka.InitProduct(golibKafak.Kafka{ServerAddr: cfg.Kafka.Address})
	if err != nil {
		panic(err)
	}
	//初始化tailf
	wg.Add(1)
	err = tailf.InitTailF(kafkaProduct, cfg, wg)
	if err != nil {
		panic(err)
	}
	wg.Wait()
}
