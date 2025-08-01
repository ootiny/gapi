package core

type TypescriptBuilder struct {
	rootConfig  GApiRootConfig
	buildConfig GApiConfig
	output      GApiRootOutputConfig
}

func (p *TypescriptBuilder) BuildServer() error {
	return nil
}

func (p *TypescriptBuilder) BuildClient() error {
	return nil
}
