package main

import (
	//"proxy_manage/utlis"
	"proxy_manage/config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"proxy_manage/pojo"

	//"proxy_manage/service"
	"proxy_manage/utlis"
	"proxy_manage/service"
	"log"
	"proxy_manage/amqphelper"
	"github.com/streadway/amqp"
	"encoding/json"
	"proxy_manage/ip2region"
	"time"
	"proxy_manage/dao"
)

func main() {
	//加载配置文件
	con := config.New("./config/config.yml")
	config.NewDb(con)
	config.GetDb().AutoMigrate(&pojo.ProxyPool{})
	ip2region.InitRegion("./config/ip2region.db")
	go utlis.Runtime(3, func() {
		log.Println("----开启代理同步任务-----")
		service.Callback(con)
	})
	go utlis.Runtime(3, func() {
		log.Println("---开启代理校验任务----")
		service.Check(con.ValidTime)
	})
	amqpConfig := &amqphelper.Config{
		AmqpUrl: con.AmqUrl,
		FromQueue: amqphelper.FromQueue{
			Name: con.WorkerName,
		},
	}
	mq := amqphelper.New(amqpConfig)
	mq.StartAsyncConsume(func(buff []byte, delivery amqp.Delivery) {
		//业务处理
		domain := pojo.Domain{}
		err := json.Unmarshal(buff, &domain)
		if err != nil {
			log.Fatal(err)
		}
		domain.IsCn = ip2region.IsCn(domain.Domain)
		domain.CreatedAt = time.Now().Unix()
		dao.SvaeDomain(domain)
		delivery.Ack(false)
	})


}
