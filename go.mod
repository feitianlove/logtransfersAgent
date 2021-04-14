module github.com/feitianlove/logtransfersAgent

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.27.2
	github.com/feitianlove/golib v0.0.0-20210412065010-e969bef46491
	github.com/hpcloud/tail v1.0.0
	github.com/sirupsen/logrus v1.8.1
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
