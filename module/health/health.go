package health

import (
	"sync"

	"gin_template/server"
	"gin_template/utils/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	instance = &health{}
	server.RegisterModule(instance)
}

var instance *health

type health struct {
	log *logrus.Entry
}

func (m *health) GetModuleInfo() server.ModuleInfo {
	return server.ModuleInfo{
		ID:       "module.health",
		Instance: instance,
	}
}

func (m *health) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
	m.log = logrus.WithField("module", "health")
}

func (m *health) PostInit(server *server.Server) {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *health) Serve(server *server.Server) {
	// 注册服务函数部分
	server.HttpEngine.GET("/health", healthHandler(m.log.WithField("func", "healthHandler"), server))
}

func (m *health) Start(server *server.Server) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *health) Stop(server *server.Server, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

var (
	_ = types.RouterGenerator(healthHandler)
)

func healthHandler(log *logrus.Entry, server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = server.DataSource.MainDB
		c.JSON(200, gin.H{
			"msg":        "pong",
			"User-Agent": c.GetHeader("User-Agent"),
		})
		instance.log.Println("Data")
		return
	}
}
