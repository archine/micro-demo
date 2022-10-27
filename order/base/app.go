package base

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.avatarworks.com/servers/component/hj-gin/plugin"
	ioc "gitlab.avatarworks.com/servers/component/hj-ioc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"micro-demo/api/user"
	"micro-demo/order/base/discovery"
	"micro-demo/order/config"
	"micro-demo/user/util"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Engine     *gin.Engine
	GrpcServer *grpc.Server
}

// NewApp 初始化应用程序
func NewApp(configPath string) (*App, error) {
	config.InitConfig(configPath)
	engine := initGin()
	initGrpcClient()
	return &App{
		Engine:     engine,
		GrpcServer: grpc.NewServer(),
	}, nil
}

// Run 启动应用程序
func (a *App) Run() {
	ipAddr, err := util.GetIpAddr()
	if err != nil {
		log.Fatalf("get local ip faild,%s", err.Error())
	}
	svc := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Port),
		Handler: a.Engine,
	}

	registrar := discovery.NewRegistrar(config.Conf.Etcd.Addr, config.Conf.AppName, config.Conf.Etcd.HeartBeat)
	registrar.Register(fmt.Sprintf("%s:%d", ipAddr, config.Conf.GrpcPort))

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case <-ch:
			registrar.Deregister()
			os.Exit(0)
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ipAddr, config.Conf.GrpcPort))
		if err != nil {
			log.Fatalf("initializer grpc server failed,%s", err.Error())
		}
		err = a.GrpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("listen to grpc server failed,%s", err.Error())
		}
	}()

	log.Infof("%s successfully started on Ports:[%d,%d]", config.Conf.AppName, config.Conf.Port, config.Conf.GrpcPort)
	if err := svc.ListenAndServe(); err != nil {
		registrar.Deregister()
		log.Fatalf("%s failed to start, %s", config.Conf.AppName, err.Error())
	}
}

// 加载grpc客户端服务
func initGrpcClient() {
	builder := discovery.NewBuilder(config.Conf.Etcd.Addr)
	resolver.Register(builder)
	dial, err := grpc.Dial(builder.Scheme()+":/"+"user", grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("init user service client failed,%s", err.Error())
	}
	ioc.SetBeans(user.NewUserClient(dial))
}

// 初始化gin
func initGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	plugin.InitLog(config.Conf.LogLevel)
	engine.Use(plugin.LogMiddleware())
	engine.Use(plugin.GlobalExceptionInterceptor)
	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))
	engine.MaxMultipartMemory = config.Conf.MaxFileSize
	engine.RemoveExtraSlash = true
	log.Debugf("init gin engine successful...")
	return engine
}
