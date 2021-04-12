package tailf

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/feitianlove/golib/kafka"
	"github.com/feitianlove/golib/tailf"
	"github.com/feitianlove/logtransfersAgent/config"
	"github.com/feitianlove/logtransfersAgent/lib"
	"github.com/feitianlove/logtransfersAgent/logger"
	"github.com/sirupsen/logrus"
	"sync"
)

var tailObjMgr *tailf.TailObjMgr

func init() {
	tailObjMgr = &tailf.TailObjMgr{
		MsgChan: make(chan *tailf.TextMsg, 10000),
	}
}

func InitTailF(product *kafka.KProduct, cfg *config.Config, wg sync.WaitGroup) error {
	fileName, err := lib.RecursionGetAllFile(cfg.TailF.Dir)
	logger.Ctrl.WithFields(logrus.Fields{
		"cfg.dir": cfg.TailF.Dir,
		"AllFIle": fileName[:10],
	}).Info("get All file")
	if err != nil {
		return err
	}
	RunTailF(fileName[:10], product)
	wg.Done()
	return nil
}

func RunTailF(fileName []string, produce *kafka.KProduct) {
	go func() {
		err := tailf.CreateTailFInstance(fileName, tailObjMgr)
		logger.Ctrl.WithFields(logrus.Fields{
			"err":      err,
			"fileName": fileName,
		}).Info("CreateTailFInstance")

		if err != nil {
			panic(err)
		}
	}()
	for i := 0; i < 10; i++ {
		go func() {
			for {
				msg := tailf.GetOneLine(tailObjMgr)
				//fmt.Println(msg)
				m := fmt.Sprintf("%s", msg)
				err, partition, offset := produce.SendMessage(&sarama.ProducerMessage{
					Topic: "ftfeng-test3",
					Key:   nil,
					Value: sarama.StringEncoder(m),
					//Headers:   nil,
					//Metadata:  nil,
					//Offset:    0,
					//Partition: 0,
					//Timestamp: time.Time{},
				})
				fmt.Println(err)
				if err != nil {
					logger.Console.WithFields(logrus.Fields{
						"Value":     m,
						"partition": partition,
						"offset":    offset,
						"fileName":  fileName,
						"err":       err,
					}).Error("kafka send message")
				} else {
					logger.Console.WithFields(logrus.Fields{
						"Value":     m,
						"partition": partition,
						"offset":    offset,
						"fileName":  fileName,
					}).Info("kafka send message")
				}
			}
		}()
	}
}
