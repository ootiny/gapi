package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GolangBuilder struct {
	rootConfig  GApiRootConfig
	buildConfig GApiConfig
	output      GApiRootOutputConfig
}

func (p *GolangBuilder) BuildServer() error {
	if p.buildConfig.Package == "" {
		return fmt.Errorf("package is required")
	}

	outDir := filepath.Join(p.output.Dir, p.buildConfig.Package)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	header := fmt.Sprintf(`// %s: %s
package %s
`, BuilderStartTag, BuilderDescription, p.buildConfig.Package)

	imports := []string{}

	defines := []string{}

	actions := []string{}

	for name, define := range p.buildConfig.Definitions {
		if define.Import != nil {
			if len(define.Attributes) > 0 {
				return fmt.Errorf("%s can not set attributes when imported", name)
			}

		}

		// defineContent := fmt.Sprintf(`type %s struct {

		// }`, define.Name)
	}

	content := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n//%s",
		header,
		strings.Join(imports, "\n"),
		strings.Join(defines, "\n"),
		strings.Join(actions, "\n"),
		BuilderEndTag,
	)

	return os.WriteFile(filepath.Join(outDir, "gapi.go"), []byte(content), 0644)
}

func (p *GolangBuilder) BuildClient() error {
	return nil
}
