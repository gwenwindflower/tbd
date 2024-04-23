package internal

import (
	"os"
	"testing"

	"github.com/gwenwindflower/tbd/shared"
)

func TestWriteStagingModels(t *testing.T) {
	ts := shared.SourceTables{
		SourceTables: []shared.SourceTable{
			{
				Name:   "table1",
				Schema: "raw",
				Columns: []shared.Column{
					{
						Name:     "COLUMN3",
						DataType: "numbers",
					},
				},
				DataTypeGroups: map[string][]shared.Column{
					"text": {
						shared.Column{
							Name:     "column1",
							DataType: "text",
						},
						shared.Column{
							Name:     "column2",
							DataType: "text",
						},
					},
					"numbers": {
						shared.Column{
							Name:     "COLUMN3",
							DataType: "numbers",
						},
					},
				},
			},
		},
	}
	expect := `with

source as (

    select * from {{ source('raw', 'table1') }}

),

renamed as (

    select
        -- numbers
        COLUMN3 as column3,
        
        -- text
        column1 as column1,
        column2 as column2,
        
    from source
)

select * from renamed
`
	bd := t.TempDir()
	WriteStagingModels(ts, bd, "stg")
	got, err := os.ReadFile(bd + "/stg_table1.sql")
	if err != nil {
		t.Errorf("Error reading staging file %v", err)
	}
	if string(got) != expect {
		t.Errorf("Expected %v, got %v", expect, string(got))
	}
}
