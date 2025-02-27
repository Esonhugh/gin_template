package __APPNAME__

import (
	"gin_template/server"
	s "gin_template/server"
	"gin_template/utils/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Register__ROUTER__(m *__APPNAME__, server *server.Server) {
	server.HttpEngine.GET("/__ROUTER__", __ROUTER__(m.log.WithField("func", "__ROUTER__"), server))
}

var (
	_ = types.RouterGenerator(__ROUTER__)
)

func __ROUTER__(l *logrus.Entry, server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := s.CreateTraceLogger(l, c)
		_ = server.DataSource.MainDB
		c.JSON(200, gin.H{
			"msg":        "pong",
			"User-Agent": c.GetHeader("User-Agent"),
		})
		log.Info("health check")
		return
	}
}
