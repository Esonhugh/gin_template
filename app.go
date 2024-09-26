package main

import (
	"gin_template/cmd"
	_ "gin_template/cmd/newroute"
	_ "gin_template/cmd/serve"

	_ "gin_template/cmd/createapp"
)

func main() {
	cmd.Execute()
}
