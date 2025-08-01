package core

import "fmt"

type GolangBuilder struct {
	rootConfig  GApiRootConfig
	buildConfig GApiConfig
}

func (p *GolangBuilder) BuildServer() error {
	if p.buildConfig.Package == "" {
		return fmt.Errorf("package is required")
	}

	content := fmt.Sprintf(`// %s: %s
package %s
`, BuilderTag, BuilderDescription, p.buildConfig.Package)

	fmt.Println(content)

	return nil
}

func (p *GolangBuilder) BuildClient() error {
	return nil
}
