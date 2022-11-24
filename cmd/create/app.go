package create

import (
	"bytes"
	"errors"
	"fmt"
	"gin_template"
	fs "gin_template/utils"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	appName string
	force   bool
)

var StartCmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new app",
	Example: "app create -n users",
	Run: func(cmd *cobra.Command, args []string) {
		err := load()
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&appName, "name", "n", "", "create a new app with provided name")
	StartCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force generate the app")
}

func load() error {
	if appName == "" {
		return errors.New("app name should not be empty, use -n")
	}

	m := map[string]string{}
	m["appNameExport"] = strings.ToUpper(appName[:1]) + appName[1:]
	m["appName"] = strings.ToLower(appName[:1]) + appName[1:]

	dir := path.Join("module", appName)

	_ = fs.IsNotExistMkDir(dir)

	if !force && (fs.FileExist(path.Join(dir, m["appName"]+".go")) || fs.FileExist(path.Join(dir, "service_test.go"))) {
		return errors.New("target file already exist, use -f flag to cover")
	}

	service := path.Join(dir, m["appName"]+".go")

	if rt, err := template.ParseFiles("template/service.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		_ = rt.Execute(&b, m)
		fs.FileCreate(b, service)
	}
	str := string(fs.ReadFile("main/server_main.go"))
	str = strings.Replace(str, "\t// New Service Add There [No Delete]\n",
		fmt.Sprintf("\t_ \"%v/module/%v\"\n\t// New Service Add There [No Delete]\n",
			gin_template.PackageName(), m["appName"]), 1)
	fs.FileCreate(*bytes.NewBufferString(str), "main/server_main.go")

	println("App " + appName + " generate success under " + dir)
	return nil
}
