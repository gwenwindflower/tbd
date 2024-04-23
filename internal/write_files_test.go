package internal

import (
	"strings"
	"testing"

	"github.com/gwenwindflower/tbd/shared"
)

func TestWriteFiles(t *testing.T) {
	ts := shared.SourceTables{
		SourceTables: []shared.SourceTable{
			{
				Name: "table1",
				Columns: []shared.Column{
					{
						Name:     "column1",
						DataType: "type1",
					},
				},
			},
		},
	}
	bd := t.TempDir()
	WriteFiles(ts, bd, "prefix")
}

func TestWriteFilesError(t *testing.T) {
	ts := shared.SourceTables{
		SourceTables: []shared.SourceTable{},
	}
	bd := t.TempDir()

	err := WriteFiles(ts, bd, "prefix")
	if err == nil {
		t.Error("expected error, got nil")
	} else {
		if !strings.Contains(err.Error(), "no tables to write") {
			t.Errorf("expected error to contain 'no tables to write', got %v", err)
		}
	}
}
