package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type GApiConfig struct {
	Listen  string   `json:"listen"`
	Output  []string `json:"output"`
	Project string   `json:"project"`
}

type IBuilder interface {
	BuildImport() (string, error)
	BuildClass() (string, error)
	BuildServerAction() (string, error)
	BuildClientAction() (string, error)
}

func LoadConfig() (GApiConfig, string, error) {
	configPath := ""
	configContent := ""

	if len(os.Args) > 1 {
		if _, err := os.Stat(os.Args[1]); err == nil {
			content, err := os.ReadFile(os.Args[1])
			if err == nil {
				configContent = string(content)
				configPath = os.Args[1]
			}
		}
	}

	if configContent == "" {
		// 在当前目录下，依次寻找 .gapi.json .gapi.yaml .gapi.yml
		searchFiles := []string{"./.gapi.json", "./.gapi.yaml", "./.gapi.yml"}
		for _, file := range searchFiles {
			if _, err := os.Stat(file); err == nil {
				content, err := os.ReadFile(file)
				if err == nil {
					configContent = string(content)
					configPath = file
					break
				}
			}
		}
	}

	if !filepath.IsAbs(configPath) {
		if absPath, err := filepath.Abs(configPath); err != nil {
			return GApiConfig{}, "", fmt.Errorf("failed to convert config path to absolute path: %v", err)
		} else {
			configPath = absPath
		}
	}

	var config GApiConfig

	switch filepath.Ext(configPath) {
	case ".json":
		err := json.Unmarshal([]byte(configContent), &config)
		if err != nil {
			log.Fatalf("Failed to parse config: %v", err)
		}
	case ".yaml", ".yml":
		err := yaml.Unmarshal([]byte(configContent), &config)
		if err != nil {
			log.Fatalf("Failed to parse config: %v", err)
		}
	default:
		return GApiConfig{}, "", fmt.Errorf("unsupported config file extension: %s", filepath.Ext(configPath))
	}

	if !filepath.IsAbs(config.Project) {
		configDir := filepath.Dir(configPath)
		config.Project = filepath.Join(configDir, config.Project)
	}

	return config, configPath, nil
}

func Output(config GApiConfig, builder IBuilder) error {
	log.Printf("Start build gapi\n")
	log.Printf("Project: %s\n", config.Project)

	return nil
}
