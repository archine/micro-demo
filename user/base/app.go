package base

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.avatarworks.com/servers/component/hj-gin/plugin"
	"google.golang.org/grpc"
	"micro-demo/user/base/discovery"
	"micro-demo/user/config"
	"micro-demo/user/util"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type App struct {
	Engine     *gin.Engine
	GrpcServer *grpc.Server
}

// NewApp 初始化应用程序
func NewApp(configPath string) (*App, error) {
	config.InitConfig(configPath)
	initPlugin()
	grpcServer := grpc.NewServer()
	engine := initGin(grpcServer)
	return &App{
		Engine:     engine,
		GrpcServer: grpcServer,
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
	registrar.Register(ipAddr)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case <-ch:
			registrar.Deregister()
			os.Exit(0)
		}
	}()

	log.Infof("%s successfully started on Ports:[%d]", config.Conf.AppName, config.Conf.Port)
	if err := svc.ListenAndServe(); err != nil {
		registrar.Deregister()
		log.Fatalf("%s failed to start, %s", config.Conf.AppName, err.Error())
	}
}

// 初始化插件
func initPlugin() {
	plugin.InitLog(config.Conf.LogLevel)
	log.Debugf("init plugin successful...")
}

// 初始化gin
func initGin(grpcServer *grpc.Server) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(plugin.LogMiddleware())
	engine.Use(plugin.GlobalExceptionInterceptor)
	engine.Use(func(context *gin.Context) {
		if context.Request.ProtoMinor == 2 && strings.HasPrefix(context.GetHeader("Content-Type"), "application/grpc") {
			context.Status(http.StatusOK)
			grpcServer.ServeHTTP(context.Writer, context.Request)
			context.Abort()
			return
		}
		context.Next()
	})
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
