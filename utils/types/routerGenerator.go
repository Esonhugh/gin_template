package types

import (
	"gin_template/server"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GinRouter func(
	log *logrus.Entry,
	server *server.Server,
) gin.HandlerFunc
