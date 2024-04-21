package main

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestInferColumnFields(t *testing.T) {
	ts := CreateTempSourceTables()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.groq.com/openai/v1/chat/completions",
		httpmock.NewStringResponder(200, `{"choices": [{"index": 0, "message": {"role": "assistant","content": "lord of rivendell"}}]}`))
	llm, err := GetLlm(FormResponse{Llm: "groq"})
	if err != nil {
		t.Errorf("Did not expect err getting LLM: %v", err)
	}
	g, ok := llm.(*Groq)
	if !ok {
		t.Error("Expceted Groq LLM type")
	}
	err = InferColumnFields(g, ts)
	if err != nil {
		t.Errorf("Did not expect err infering column fields: %v", err)
	}

	info := httpmock.GetCallCountInfo()
	if info["POST https://api.groq.com/openai/v1/chat/completions"] != 2 {
		t.Error("expected", 2, "got", info["POST https://api.groq.com/openai/v1/chat/completions"])
	}

	expected := "lord of rivendell"
	desc := ts.SourceTables[0].Columns[0].Description
	if desc != expected {
		t.Error("expected", expected, "got", desc)
	}
}
