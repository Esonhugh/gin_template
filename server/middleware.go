package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var ginLogger = logrus.WithField("server", "gin")

func GetTraceID(c *gin.Context) string {
	traceID, _ := c.Get("request_trace_id")
	return traceID.(string)
}

func CreateTraceLogger(log *logrus.Entry, c *gin.Context) *logrus.Entry {
	traceID := GetTraceID(c)
	return log.WithField("trace_id", traceID)
}
func TraceRequest(c *gin.Context) {
	uid := uuid.New()
	c.Set("request_trace_id", uid.String())
}

func ginRequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		ginLogger.Infof("| %3d | %13v | %15s | %s | %s | trace_id=%s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			GetTraceID(c),
		)
	}
}
