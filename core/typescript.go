package core

type TypescriptBuilder struct {
	output GApiOutputConfig
	config GApiConfig
}

func (p *TypescriptBuilder) BuildServer() (string, error) {
	return "", nil
}

func (p *TypescriptBuilder) BuildClient() (string, error) {
	return "", nil
}
