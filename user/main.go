package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.avatarworks.com/servers/component/hj-gin/mvc"
	ioc "gitlab.avatarworks.com/servers/component/hj-ioc"
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
