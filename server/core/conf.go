package core

import (
	"fmt"
	"log"
	"server/config"
	"server/utils"

	"gopkg.in/yaml.v3"
)

func InitConf() *config.Config {
	c := &config.Config{}
	yamlConf, err := utils.LoadYAML()
	fmt.Println("-------", string(yamlConf), "-------------")
	if err != nil {
		log.Fatalf("Failed to load configuration :%v", err)
	}
	err = yaml.Unmarshal(yamlConf, c)
	fmt.Println("-------", c, "-------------")
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML configuration:%v", err)
	}
	return c
}
