package newroute

import (
	"bytes"
	_ "embed"
	"fmt"
	"gin_template/utils/file"
	"strings"

	"gin_template/cmd"
	"github.com/spf13/cobra"
)

var appName string

func init() {
	cmd.RootCmd.AddCommand(RouterCmd)
	RouterCmd.PersistentFlags().StringVarP(&appName, "app", "a", "", "app name")
}

var RouterCmd = &cobra.Command{
	Use:   "route",
	Short: "Create new router function ",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			println("router name should not be empty! ")
			return
		}
		if appName != "" {
			_ = file.IsNotExistMkDir("module/" + appName)
		}
		Template := string(file.ReadFile("template/router.go"))
		for _, router := range args {
			routerContent := strings.ReplaceAll(
				Template,
				"__ROUTER__",
				router)
			if appName != "" {
				routerContent = strings.ReplaceAll(
					routerContent,
					"__APPNAME__",
					appName)
				fn := "module/" + appName + "/" + router + ".go"
				if file.FileExist(fn) {
					println("target file '" + fn + "' already exist, skipped")
					continue
				}
				var b = bytes.NewBufferString(routerContent)
				file.FileCreate(*b, fn)
				println("create file '" + fn + "' success")
			} else {
				fmt.Println(routerContent)
			}
		}
	},
}
