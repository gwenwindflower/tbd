package shared

type Column struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	DataType    string   `yaml:"data_type"`
	Tests       []string `yaml:"tests"`
}

type SourceTable struct {
	DataTypeGroups map[string][]Column `yaml:"-"`
	Name           string              `yaml:"name"`
	Schema         string              `yaml:"-"`
	Columns        []Column            `yaml:"columns"`
}

type SourceTables struct {
	SourceTables []SourceTable `yaml:"sources"`
}

type ConnectionDetails struct {
	ProjectName    string
	ConnType       string
	Account        string
	Database       string
	Dataset        string
	Project        string
	Username       string
	Path           string
	Schema         string
	Host           string
	Password       string
	PasswordEnvVar string
	SslMode        string
	Catalog        string
	Token          string
	TokenEnvVar    string
	HttpPath       string
	Port           int
}
