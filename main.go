package main

import (
	"BlumBot/core"
	"BlumBot/tools"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("configs/config.yml")
	if err != nil {
		panic(err)
	}

	tools.PrintLogo()

	core.LaunchBot()
}
