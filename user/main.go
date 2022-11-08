package main

import (
	"github.com/archine/gin-plus/mvc"
	"github.com/archine/ioc"
	log "github.com/sirupsen/logrus"
	"micro-demo/user/base"
	_ "micro-demo/user/controller"
)

func main() {
	app, err := base.NewApp("user/app.yml")
	if err != nil {
		log.Fatalln(err.Error())
	}
	ioc.SetBeans(app.GrpcServer)
	mvc.Apply(app.Engine, true)
	app.Run()
}
