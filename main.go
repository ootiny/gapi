package main

import (
	"log"

	"github.com/ootiny/gapi/core"
)

func main() {
	config, configPath, err := core.LoadConfig()
	if err != nil {
		log.Panicf("Failed to load config: %v", err)
	}

	log.Printf("using config file: %s", configPath)

	for _, output := range config.Output {
		switch output {
		case "golang":
			core.Output(config, &core.GolangBuilder{})
		case "typescript":
			core.Output(config, &core.TypescriptBuilder{})
		default:
			log.Panicf("Unsupported output: %s", output)
		}
	}
}
