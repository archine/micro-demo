package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.avatarworks.com/servers/component/hj-gin/mvc"
	ioc "gitlab.avatarworks.com/servers/component/hj-ioc"
	"micro-demo/order/base"
	_ "micro-demo/order/controller"
)

func main() {
	app, err := base.NewApp("order/app.yml")
	if err != nil {
		log.Fatalln(err.Error())
	}
	ioc.SetBeans(app.GrpcServer)
	mvc.Apply(app.Engine, true)
	app.Run()
}
