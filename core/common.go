package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

type GApiRootConfig struct {
	Listen  string `json:"listen"`
	Project string `json:"project"`
}

type gApiConfigHeader struct {
	Version string `json:"version"`
}

type IBuilder interface {
	BuildImport() (string, error)
	BuildClass() (string, error)
	BuildServerAction() (string, error)
	BuildClientAction() (string, error)
}

func LoadRootConfig() (GApiRootConfig, string, error) {
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
			return GApiRootConfig{}, "", fmt.Errorf("failed to convert config path to absolute path: %v", err)
		} else {
			configPath = absPath
		}
	}

	var config GApiRootConfig

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
		return GApiRootConfig{}, "", fmt.Errorf("unsupported config file extension: %s", filepath.Ext(configPath))
	}

	if !filepath.IsAbs(config.Project) {
		configDir := filepath.Dir(configPath)
		config.Project = filepath.Join(configDir, config.Project)
	}

	return config, configPath, nil
}

func Output(config GApiRootConfig) error {
	log.Printf("Start build gapi\n")
	log.Printf("Project Dir: %s\n", config.Project)

	versions := []string{"gapi", "gapi.v1"}

	walkErr := filepath.Walk(config.Project, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("warn: failed to read file %s: %v", path, err)
				return nil // continue walking
			}

			var header gApiConfigHeader
			if err := json.Unmarshal(content, &header); err != nil {
				return nil // Not a gapi config file, just ignore.  continue walking
			}

			if slices.Contains(versions, header.Version) {
				if err := OutputFile(path); err != nil {
					return err // stop walking and return error
				}
			}
		}
		return nil
	})

	if walkErr != nil {
		return fmt.Errorf("error walking project directory: %w", walkErr)
	}

	return nil
}

func OutputFile(absPath string) error {
	fmt.Println(absPath)

	return nil
}
