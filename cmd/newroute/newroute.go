package newroute

import (
	"fmt"
	"strings"

	"gin_template/server"
	"gin_template/utils/types"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Template = `
var (
	_ = types.RouterGenerator(__ROUTER__)
)

func __ROUTER__ (log *logrus.Entry, server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = server.DataSource.MainDB 
	}
}
`

func TestRouter(log *logrus.Entry, server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = server.DataSource.MainDB
	}
}

var (
	_      = types.RouterGenerator(TestRouter)
	router string
)

func init() {
	RouterCmd.PersistentFlags().StringVarP(
		&router, "router", "n", "ping", "Router handler generate functions",
	)
}

var RouterCmd = &cobra.Command{
	Use:   "route",
	Short: "Create new router function ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.ReplaceAll(
			Template,
			"__ROUTER__",
			router))
	},
}
