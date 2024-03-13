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
	Warehouse string
	Username  string
	Account   string
	Schema    string
	Database  string
}
