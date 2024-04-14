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
	Columns        []Column            `yaml:"columns"`
}

type SourceTables struct {
	SourceTables []SourceTable `yaml:"sources"`
}

type ConnectionDetails struct {
	ConnType string
	Username string
	Account  string
	Database string
	Schema   string
	Project  string
	Dataset  string
	Path     string
}
