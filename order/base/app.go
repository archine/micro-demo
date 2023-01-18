package base

import (
	"fmt"
	"github.com/archine/gin-plus/v2/ast"
	"github.com/archine/gin-plus/v2/exception"
	"github.com/archine/gin-plus/v2/mvc"
	"github.com/archine/gin-plus/v2/plugin"
	"github.com/archine/ioc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

// RunApplication run the main program entry
// @Param args start command parameter
func RunApplication(args []string, globalFunc ...gin.HandlerFunc) {
	if len(args) > 1 && os.Args[1] == "ast" {
		ast.Parse()
		return
	}
	config.LoadApplicationConfigFile("order/app.yml")
	startApplication(initGin(), globalFunc, initGrpcClient())
}

// 加载grpc服务
func initGrpcClient() *grpc.Server {
	builder := discovery.NewBuilder(config.Conf.Etcd.Addr)
	resolver.Register(builder)
	dial, err := grpc.Dial(builder.Scheme()+":/"+"user", grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("init user service client failed,%s", err.Error())
	}
	ioc.SetBeans()
	grpcServer := grpc.NewServer()
	ioc.SetBeans(user.NewUserClient(dial), grpcServer)
	return grpcServer
}

// 初始化gin
func initGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	plugin.InitLog(config.Conf.LogLevel)
	engine.Use(plugin.LogMiddleware())
	engine.Use(exception.GlobalExceptionInterceptor)
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
	ioc.SetBeans(engine)
	log.Debugf("init gin engine successful...")
	return engine
}

func startApplication(engine *gin.Engine, globalFunc []gin.HandlerFunc, grpcServer *grpc.Server) {
	mvc.Apply(engine, true, Ast, globalFunc...)
	log.Debugf("api apply successful...")
	ipAddr, err := util.GetIpAddr()
	if err != nil {
		log.Fatalf("get local ip faild,%s", err.Error())
	}
	svc := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Port),
		Handler: engine,
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
		var listener net.Listener
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", ipAddr, config.Conf.GrpcPort))
		if err != nil {
			log.Fatalf("initializer grpc server failed,%s", err.Error())
		}
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("listen to grpc server failed,%s", err.Error())
		}
	}()

	log.Infof("%s successfully started on Ports:[%d,%d]", config.Conf.AppName, config.Conf.Port, config.Conf.GrpcPort)
	if err = svc.ListenAndServe(); err != nil {
		registrar.Deregister()
		log.Fatalf("%s failed to start, %s", config.Conf.AppName, err.Error())
	}
}
