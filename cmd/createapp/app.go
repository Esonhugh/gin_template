package createapp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"gin_template/cmd"
	"gin_template/utils"
	"gin_template/utils/file"
	"github.com/spf13/cobra"
)

var (
	force bool
)

var StartCmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new app",
	Example: "app create users",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 0 {
			println("app name should not be empty! ")
			os.Exit(1)
		}
		for _, appName := range args {
			err := load(appName)
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
		}
	},
}

func init() {
	StartCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force generate the app")
	cmd.RootCmd.AddCommand(StartCmd)
}

func load(appName string) error {
	if appName == "" {
		return errors.New("app name should not be empty, use -n")
	}

	var m = make(map[string]string)
	dir := path.Join("module", appName)
	{
		_ = file.IsNotExistMkDir(dir)

		if !force && (file.FileExist(path.Join(dir, appName+".go")) || file.FileExist(path.Join(dir, "service_test.go"))) {
			return errors.New("target file already exist, use -f flag to cover")
		}
		service := path.Join(dir, appName+".go")

		b := file.ReadFile("template/service.go")
		tmpl, err := template.New("service").Parse(strings.ReplaceAll(string(b), "__APPNAME__", appName))
		if err != nil {
			return errors.New("parse template error: " + err.Error())
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, m); err != nil {
			return errors.New("execute template error: " + err.Error())
		} else {
			file.FileCreate(buf, service)
		}
	}
	{
		serverMain := path.Join(
			"cmd", "serve", "server.go")
		str := string(file.ReadFile(serverMain))
		str = strings.Replace(str, "\t// New Service Add There [No Delete]\n",
			fmt.Sprintf("\t_ \"%v/module/%v\"\n\t// New Service Add There [No Delete]\n",
				utils.PackageName(), appName), 1)
		file.FileCreate(*bytes.NewBufferString(str), serverMain)
	}
	println("App " + appName + " generate success under " + dir)
	return nil
}
