package main

import (
	"log"

	"github.com/ootiny/gapi/core"
)

func main() {
	rootConfig, configPath, err := core.LoadRootConfig()
	if err != nil {
		log.Panicf("Failed to load config: %v", err)
	}

	log.Printf("using config file: %s", configPath)
	if err := core.Output(rootConfig); err != nil {
		log.Panicf("Failed to output: %v", err)
	}
}
