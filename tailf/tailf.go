package tailf

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/feitianlove/golib/kafka"
	"github.com/feitianlove/golib/tailf"
	"github.com/feitianlove/logtransfersAgent/config"
	"github.com/feitianlove/logtransfersAgent/lib"
	"github.com/feitianlove/logtransfersAgent/logger"
	"github.com/sirupsen/logrus"
	"regexp"
	"sync"
	"time"
)

var (
	tailObjMgr  *tailf.TailObjMgr
	FileNameMap map[string]string
	Lock        *sync.RWMutex
)

func init() {
	tailObjMgr = &tailf.TailObjMgr{
		MsgChan: make(chan *tailf.TextMsg, 10000),
	}
	FileNameMap = make(map[string]string, 0)
	Lock = new(sync.RWMutex)
}

func InitTailF(product *kafka.KProduct, cfg *config.Config, wg sync.WaitGroup) error {
	fileName, err := lib.RecursionGetAllFile(cfg.TailF.Dir)
	if err != nil {
		return err
	}
	RunTailF(cfg, fileName, product)
	wg.Done()
	return nil
}

func RunTailF(cfg *config.Config, fileName []string, product *kafka.KProduct) {
	CreateTailF(fileName)
	GetDataFromTailChan(product, cfg)
	scanfNewFileTime := time.NewTicker(60 * time.Second)
	expireFileMapTime := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			CreateUpdateTailFile(cfg)
			<-scanfNewFileTime.C
		}
	}()
	go func() {
		for {
			LeastRecently()
			<-expireFileMapTime.C
		}
	}()
}

func CreateTailF(fileName []string) {
	go func() {
		Lock.Lock()
		defer Lock.Unlock()
		for _, item := range fileName {
			err := tailf.CreateTailFInstance([]string{item}, tailObjMgr)
			if err != nil {
				logger.Ctrl.WithFields(logrus.Fields{
					"err":      err,
					"fileName": item,
				}).Error("CreateTailFInstance")
			} else {
				logger.Ctrl.WithFields(logrus.Fields{
					"err":      err,
					"fileName": item,
				}).Info("CreateTailFInstance")
			}

			if _, ok := FileNameMap[item]; !ok && err == nil {
				FileNameMap[item] = "1"
			}
		}
	}()
}

func CreateUpdateTailFile(cfg *config.Config) {
	fileName := []string{}
	files, err := lib.RecursionGetAllFile(cfg.TailF.Dir)
	if err != nil {

	}
	Lock.RLock()
	defer Lock.RUnlock()
	for _, item := range files {
		if _, ok := FileNameMap[item]; !ok {
			fileName = append(fileName, item)
		}
	}
	if len(fileName) == 0 {
		logger.Ctrl.WithFields(logrus.Fields{
			"msg": "don't get new file",
		}).Info("CreateUpdateTailFile")
		return
	}
	CreateTailF(fileName)
}
func GetDataFromTailChan(produce *kafka.KProduct, cfg *config.Config) {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				msg := tailf.GetOneLine(tailObjMgr)
				m := fmt.Sprintf("%s", *msg)
				data, _ := json.Marshal(GetFormatData(m))
				err, partition, offset := produce.SendMessage(&sarama.ProducerMessage{
					Topic: cfg.Kafka.Topic,
					Key:   nil,
					Value: sarama.StringEncoder(string(data)),
					//Headers:   nil,
					//Metadata:  nil,
					//Offset:    0,
					//Partition: 0,
					//Timestamp: time.Time{},
				})
				if err != nil {
					logger.Console.WithFields(logrus.Fields{
						"Value":     m,
						"partition": partition,
						"offset":    offset,
						"err":       err,
					}).Error("kafka send message")
				} else {
					logger.Console.WithFields(logrus.Fields{
						"Value":     m,
						"partition": partition,
						"offset":    offset,
					}).Info("kafka send message")
				}
			}
		}()
	}
}
func GetFormatData(s string) map[string]string {
	//获取时间
	var data = make(map[string]string)
	rTime := regexp.MustCompile(`\[([^]]+?)\]`)
	Mtime := rTime.FindStringSubmatch(s)
	if len(Mtime) > 0 {
		s = s[len(Mtime[1]):]
		data["Time"] = Mtime[1]
	} else {
		return data
	}
	r := regexp.MustCompile(`([a-zA-Z0-9]*):([a-zA-Z0-9\/\[\]\-]*)`)
	// [[Op:DownloadFile Op DownloadFile]]
	subMatchs := r.FindAllStringSubmatch(s, -1)
	for _, item := range subMatchs {
		data[item[1]] = item[2]
	}
	return data
}

//  淘汰filemap中的文件名
func LeastRecently() {
	Lock.Lock()
	defer Lock.Unlock()
	timeFormat := "2006-01-02 15:04:05"
	reg := regexp.MustCompile(`log\.(.*)\.log`)
	var deleteFileList []string
	for _, fileName := range FileNameMap {

		res := reg.FindStringSubmatch(fileName)
		var fileExpireTime int64
		if len(res) > 0 && len(res[1]) == 12 {
			t := res[1]
			temp := t[:4] + "-" + t[4:6] + "-" + t[6:8] + " " + t[8:10] + ":" + t[10:12] + ":" + "00"
			fileTime, _ := time.Parse(timeFormat, temp)
			fileExpireTime = fileTime.Add(24 * time.Hour).Unix()
		} else {
			continue
		}
		lruTime := time.Now().Unix()
		if fileExpireTime < lruTime {
			deleteFileList = append(res, fileName)
			delete(FileNameMap, fileName)
		}
	}
	for _, item := range deleteFileList {
		logger.Ctrl.WithFields(logrus.Fields{
			"LeastRecently": item,
		}).Info("expire file")
	}
}
