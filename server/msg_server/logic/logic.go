package logic

import (
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qf/util/hashids"
	"github.com/qsock/qim/config/common"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/types"
	"github.com/qsock/qim/server/msg_server/config"
)

var (
	hashIdCodec *hashids.HashID

	consumers []*ka.Consumer
)

func Init() {
	initCodec()
	initConsumer()
}

func Stop() {
	for _, consumer := range consumers {
		consumer.Stop()
	}
}

func initCodec() {
	codec, err := hashids.NewWithData(&hashids.HashIDData{
		MinLength: config.GetHashId().Mini,
		Alphabet:  config.GetHashId().Characters,
		Salt:      config.GetHashId().Salt,
	})
	if err != nil {
		panic(err)
	}
	hashIdCodec = codec
}

func initConsumer() {
	{
		kacfg := &ka.Config{}
		kacfg.Brokers = types.GetKaAddrs(config.GetEnv())
		kacfg.Group = mq.ConsumserMsg
		kacfg.Topic = mq.TopicIm
		kacfg.Workers = 50
		consumer := ka.NewConsumer(kacfg, HandleImMsg)
		consumer.Run()
		consumers = append(consumers, consumer)
	}
	{
		kacfg := &ka.Config{}
		kacfg.Brokers = types.GetKaAddrs(config.GetEnv())
		kacfg.Group = mq.ConsumserMsg
		kacfg.Topic = mq.TopicEvent
		kacfg.Workers = 20
		consumer := ka.NewConsumer(kacfg, HandleServerEvent)
		consumer.Run()
		consumers = append(consumers, consumer)
	}
}

func getWsServernames() []string {
	arrs := qgrpc.GetPrefixAddrsRegisterModel(common.WsServiceName)
	keys := make([]string, 0)
	for _, v := range arrs {
		keys = append(keys, v.Name)
	}
	return keys
}
