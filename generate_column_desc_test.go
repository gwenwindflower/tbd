package main

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetGroqResponse(t *testing.T) {
	prompt := "Who destroyed Orthanc"
	httpmock.Activate()
	defer httpmock.Deactivate()
	httpmock.RegisterResponder("POST", "https://api.groq.com/openai/v1/chat/completions",
		httpmock.NewStringResponder(200, `{"choices": [{"index": 0, "message": {"role": "assistant","content": "Treebeard and the Ents destroyed Orthanc."}}]}`))
	GroqResponse, err := GetGroqResponse(prompt)
	if err != nil {
		t.Error("expected", nil, "got", err)
	}
	expected := "Treebeard and the Ents destroyed Orthanc."
	if GroqResponse.Choices[0].Message.Content != expected {
		t.Error("expected", expected, "got", GroqResponse.Choices[0].Message.Content)
	}
}

func TestGenerateColumnDescriptions(t *testing.T) {
	ts := CreateTempSourceTables()
	httpmock.Activate()
	defer httpmock.Deactivate()
	httpmock.RegisterResponder("POST", "https://api.groq.com/openai/v1/chat/completions",
		httpmock.NewStringResponder(200, `{"choices": [{"index": 0, "message": {"role": "assistant","content": "lord of rivendell"}}]}`))
	GenerateColumnDescriptions(ts)
	expected := "lord of rivendell"

	desc := ts.SourceTables[0].Columns[0].Description
	if desc != expected {
		t.Error("expected", expected, "got", desc)
	}
}
