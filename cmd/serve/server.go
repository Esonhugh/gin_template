package serve

import (
	"os"
	"os/signal"

	_ "gin_template/module/health"
	_ "gin_template/module/test"
	// New Service Add There [No Delete]
	"gin_template/cmd"
	"gin_template/server"
	"gin_template/utils/log"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(ServerCmd)
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
	// config.Init() // if you need read config
	// logrus.StandardLogger().SetLevel(logrus.TraceLevel)
	server.Init()

	server.StartService()

	server.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	server.Stop()
}
