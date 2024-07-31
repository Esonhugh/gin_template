package serve

import (
	// New Service Add There [No Delete]
	_ "gin_template/module/health"
	"gin_template/server"
	"gin_template/utils/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func init() {
	log.WriteLogToFS()
}

var ServerCmd = &cobra.Command{
	Use:   "serve",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		ServerRun()
	},
}

func ServerRun() {
	server.Init()

	server.StartService()

	server.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	server.Stop()
}
