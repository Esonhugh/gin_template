package health

import (
	"gin_template/server"
	"github.com/gin-gonic/gin"
	"sync"
)

func init() {
	instance = &health{}
	server.RegisterModule(instance)
}

var instance *health

type health struct {
}

func (m *health) GetModuleInfo() server.ModuleInfo {
	return server.ModuleInfo{
		ID:       "internal.health",
		Instance: instance,
	}
}

func (m *health) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *health) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *health) Serve(server *server.Server) {
	// 注册服务函数部分
	server.HttpEngine.GET("/health", HealthHandler)
}

func (m *health) Start(server *server.Server) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如 http 服务器等等
}

func (m *health) Stop(server *server.Server, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg":    "Server is Online",
		"status": "ok",
	})
}
