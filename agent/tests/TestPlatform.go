package main

import (
	"common/structs"
	"common/utils"
)

func main() {
	setupfilelocation := "/Users/eladh/Desktop/install-platform-go/cli/resources/setup-eplus-full.yaml"
	setupInfo := structs.SetupInfo{}
	utils.LoadYamlFile(setupfilelocation , &setupInfo)
}