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

type GApiOutputConfig struct {
	Kind     string `json:"kind" required:"true"`
	Language string `json:"language" required:"true"`
	Package  string `json:"package"`
	FilePath string `json:"filePath" required:"true"`
}

type GApiDefinitionAttributeConfig struct {
	Name        string `json:"name" required:"true"`
	Type        string `json:"type" required:"true"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type GApiDefinitionImportConfig struct {
	From string `json:"from"`
	Type string `json:"type"`
}

type GApiDefinitionConfig struct {
	Description string                          `json:"description"`
	Attributes  []GApiDefinitionAttributeConfig `json:"attributes"`
	Import      GApiDefinitionImportConfig      `json:"import"`
}

type GApiActionParameterConfig struct {
	Name        string `json:"name" required:"true"`
	Type        string `json:"type" required:"true"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type GApiActionReturnConfig struct {
	Type        string `json:"type" required:"true"`
	Description string `json:"description"`
}

type GApiActionConfig struct {
	Description string                      `json:"description"`
	Method      string                      `json:"method" required:"true"`
	Parameters  []GApiActionParameterConfig `json:"parameters"`
	Returns     []GApiActionReturnConfig    `json:"returns"`
}

type GApiConfig struct {
	Version     string                          `json:"version" required:"true"`
	ApiPath     string                          `json:"apiPath" required:"true"`
	Outputs     []GApiOutputConfig              `json:"outputs" required:"true"`
	Description string                          `json:"description"`
	Definitions map[string]GApiDefinitionConfig `json:"definitions" required:"true"`
	Actions     map[string]GApiActionConfig     `json:"actions" required:"true"`
}

type IBuilder interface {
	BuildImport() (string, error)
	BuildClass() (string, error)
	BuildServerAction() (string, error)
	BuildClientAction() (string, error)
}

func UnmarshalConfig(filePath string, v any) error {
	if content, err := os.ReadFile(filePath); err != nil {
		return err
	} else {
		switch filepath.Ext(filePath) {
		case ".json":
			return json.Unmarshal(content, v)
		case ".yaml", ".yml":
			return yaml.Unmarshal(content, v)
		default:
			return fmt.Errorf("unsupported file extension: %s", filepath.Ext(filePath))
		}
	}
}

func LoadConfig(filePath string) (GApiConfig, error) {
	var config GApiConfig

	if err := UnmarshalConfig(filePath, &config); err != nil {
		return GApiConfig{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

func LoadRootConfig() (GApiRootConfig, string, error) {
	configPath := ""

	if len(os.Args) > 1 {
		if fileInfo, err := os.Stat(os.Args[1]); err == nil && !fileInfo.IsDir() {
			configPath = os.Args[1]
		}
	}

	if configPath == "" {
		// 在当前目录下，依次寻找 .gapi.json .gapi.yaml .gapi.yml
		searchFiles := []string{"./.gapi.json", "./.gapi.yaml", "./.gapi.yml"}
		for _, file := range searchFiles {
			if fileInfo, err := os.Stat(file); err == nil && !fileInfo.IsDir() {
				configPath = file
				break
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

	if err := UnmarshalConfig(configPath, &config); err != nil {
		return GApiRootConfig{}, "", fmt.Errorf("failed to parse config file: %w", err)
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

		var header struct {
			Version string `json:"version"`
		}

		switch filepath.Ext(path) {
		case ".json", ".yaml", ".yml":
			if err := UnmarshalConfig(path, &header); err != nil {
				return nil // Not a gapi config file, just ignore.  continue walking
			} else if slices.Contains(versions, header.Version) {
				return OutputFile(path) // output file
			} else {
				return nil
			}
		default:
			return nil
		}
	})

	if walkErr != nil {
		return fmt.Errorf("error walking project directory: %w", walkErr)
	}

	return nil
}

func OutputFile(absPath string) error {

	if config, err := LoadConfig(absPath); err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	} else {
		jsonBytes, _ := json.MarshalIndent(config, "", "  ")

		fmt.Println(string(jsonBytes))

		return nil
	}

}
