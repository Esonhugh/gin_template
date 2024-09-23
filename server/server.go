package server

import (
	"sync"

	"gin_template/server/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Instance *Server

type Server struct {
	HttpEngine *gin.Engine
	DataSource database.GlobalDB
	APIRouter  *gin.RouterGroup
	// More Service like database service
}

var logger = logrus.WithField("server", "internal")

// Init 快速初始化
// Adding in init function
func Init() {
	logger.Info("Starting Register Server Instance")

	logger.Info("initializing HTTP server instance ...")
	gin.SetMode(gin.ReleaseMode)
	httpEngine := gin.New()
	httpEngine.Use(ginRequestLog(), gin.Recovery())
	logger.Info("HTTP instance complete")

	logger.Info("initializing DataSource instance ...")
	err := database.Init()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("DataSource instance complete")

	Instance = &Server{
		HttpEngine: httpEngine,
		APIRouter:  httpEngine.Group("/api"),
		DataSource: database.GlobalDatabase,
	}
}

// Run 正式开启服务
func Run() {
	go func() {
		logger.Info("http engine starting...")
		// Running Port 9955 may need Change.
		if err := Instance.HttpEngine.Run("0.0.0.0:9955"); err != nil {
			logger.Fatal(err)
		} else {
			logger.Info("http engine running...")
		}
	}()
}

// StartService 启动服务
// 根据 Module 生命周期 此过程应在Login前调用
// 请勿重复调用
func StartService() {
	logger.Infof("initializing modules ...")
	for _, mi := range modules {
		mi.Instance.Init()
	}
	for _, mi := range modules {
		mi.Instance.PostInit(Instance)
	}
	logger.Info("all modules initialized")

	logger.Info("registering modules serve functions ...")
	for _, mi := range modules {
		mi.Instance.Serve(Instance)
	}
	logger.Info("all modules serve functions registered")

	logger.Info("starting modules tasks ...")
	for _, mi := range modules {
		go mi.Instance.Start(Instance)
	}
	logger.Info("tasks running")
}

// Stop 停止所有服务
// 调用此函数并不会使Bot离线
func Stop() {
	logger.Warn("stopping ...")
	wg := sync.WaitGroup{}
	for _, mi := range modules {
		wg.Add(1)
		mi.Instance.Stop(Instance, &wg)
	}
	wg.Wait()
	logger.Info("stopped")
	modules = make(map[string]ModuleInfo)
}
