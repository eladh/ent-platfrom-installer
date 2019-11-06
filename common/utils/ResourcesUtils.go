package utils

import (
	"github.com/kris-nova/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var devMode, _ = strconv.ParseBool(os.Getenv("DEV_MODE"))

func GetCommonResources() string {
	return "common"
}

func GetAgentResources() string {
	return "agent"
}

func GetRootDir() string {
	prefix := "/home/appuser"
	if devMode {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return dir + "/.."
	}

	return prefix
}

func GetResource(resource string, resourceType string) string {
	content, err := ioutil.ReadFile(GetResourceLocation(resource, resourceType))

	if err != nil {
		logger.Critical("get resource error", err)
	}
	return string(content)
}

func IsResourceExists(resource string, resourceType string) bool {
	return FileExists(GetRootDir() + "/" + resourceType + "/resources/" + resource)
}

func GetResourceLocation(resource string, resourceType string) string {
	location := GetRootDir() + "/" + resourceType + "/resources/" + resource

	if !FileExists(location) {
		panic("file not found at location = " + location)

	}
	return location
}

func CopyResourceToLocation(resource string, targetLocation string, resourceType string) {
	dir := filepath.Dir(targetLocation)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0777)
	}

	err := ioutil.WriteFile(targetLocation, []byte(GetResource(resource, resourceType)), 0777)
	if err != nil {
		panic(err)
	}
}

func GetYamlResourceAsMap(file string, resourceType string) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	GetYamlResource(file, &m, resourceType)

	return m
}

func GetYamlResource(resource string, output interface{}, resourceType string) {
	err := yaml.Unmarshal([]byte(GetResource(resource, resourceType)), output)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
