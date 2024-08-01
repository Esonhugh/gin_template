package types

import (
	"gin_template/server"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouterGenerator func(
	log *logrus.Entry,
	server *server.Server,
) gin.HandlerFunc
