package main

import (
	"github.com/archine/gin-plus/mvc"
	"github.com/archine/ioc"
	log "github.com/sirupsen/logrus"
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
