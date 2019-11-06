package utils

import (
	"github.com/kris-nova/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadYamlFile(file string, output interface{}) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Always("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, output)
	if err != nil {
		logger.Always("Unmarshal: %v", err)
	}
}
