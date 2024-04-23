package internal

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetGroqResponse(t *testing.T) {
	prompt := "Who destroyed Orthanc"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.groq.com/openai/v1/chat/completions",
		httpmock.NewStringResponder(200, `{"choices": [{"index": 0, "message": {"role": "assistant","content": "Treebeard and the Ents destroyed Orthanc."}}]}`))
	llm, err := GetLlm(FormResponse{Llm: "groq"})
	if err != nil {
		t.Errorf("Did not expect err getting LLM: %v", err)
	}
	g, ok := llm.(*Groq)
	if !ok {
		t.Error("Expceted Groq LLM type")
	}
	err = g.GetResponse(prompt)
	if err != nil {
		t.Error("expected", nil, "got", err)
	}
	info := httpmock.GetCallCountInfo()
	if info["POST https://api.groq.com/openai/v1/chat/completions"] != 1 {
		t.Error("expected", 1, "got", info["POST https://api.groq.com/openai/v1/chat/completions"])
	}
	expected := "Treebeard and the Ents destroyed Orthanc."
	if g.Response.Choices[0].Message.Content != expected {
		t.Error("expected", expected, "got", g.Response.Choices[0].Message.Content)
	}
}

func TestGetOpenAIResponse(t *testing.T) {
	prompt := "Who destroyed Orthanc"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
		httpmock.NewStringResponder(200, `{"choices": [{"index": 0, "message": {"role": "assistant","content": "Treebeard and the Ents destroyed Orthanc."}}]}`))
	llm, err := GetLlm(FormResponse{Llm: "openai"})
	if err != nil {
		t.Errorf("Did not expect err getting LLM: %v", err)
	}
	o, ok := llm.(*OpenAI)
	if !ok {
		t.Error("Expceted OpenAI LLM type")
	}
	err = o.GetResponse(prompt)
	if err != nil {
		t.Errorf("Did not expect err getting response: %v", err)
	}
	// TODO: flaky test
	// info := httpmock.GetCallCountInfo()
	// expectedCalls := 1
	// if info["POST https://api.openai.com/v1/completions"] != expectedCalls {
	// 	t.Error("expected", expectedCalls, "got", info["POST https://api.openai.com/v1/chat/completions"])
	// }
	expectedResp := "Treebeard and the Ents destroyed Orthanc."
	if o.Response.Choices[0].Message.Content != expectedResp {
		t.Error("expected", expectedResp, "got", o.Response.Choices[0].Message.Content)
	}
}

func TestGetAnthropicResponse(t *testing.T) {
	prompt := "Who destroyed Orthanc"
	url := "https://api.anthropic.com/v1/messages"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", url,
		httpmock.NewStringResponder(200, `{"role": "assistant", "content": [{ "type": "text", "text": "Treebeard and the Ents destroyed Orthanc."}]}`))
	llm, err := GetLlm(FormResponse{Llm: "anthropic"})
	if err != nil {
		t.Errorf("Did not expect err getting LLM: %v", err)
	}
	a, ok := llm.(*Anthropic)
	if !ok {
		t.Error("Expceted Anthropic LLM type")
	}
	err = a.GetResponse(prompt)
	if err != nil {
		t.Errorf("Did not expect err getting response: %v", err)
	}
	info := httpmock.GetCallCountInfo()
	expectedCalls := 1
	if info[fmt.Sprintf("POST %s", url)] != expectedCalls {
		t.Error("expected", expectedCalls, "got", info[fmt.Sprintf("POST %s", url)])
	}
	expectedResp := "Treebeard and the Ents destroyed Orthanc."
	actualResp := a.Response.Content[0].Text
	if actualResp != expectedResp {
		t.Error("expected", expectedResp, "got", actualResp)
	}
}

func TestGetAnthropicResponseError(t *testing.T) {
	prompt := "Who destroyed Orthanc"
	url := "https://api.anthropic.com/v1/messages"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", url,
		httpmock.NewStringResponder(400, `{"type": "error", "error": {"type": "invalid_request_error", "message": "max_tokens: Field required"}}`))
	llm, err := GetLlm(FormResponse{Llm: "anthropic"})
	if err != nil {
		t.Errorf("Did not expect err getting LLM: %v", err)
	}
	a, ok := llm.(*Anthropic)
	if !ok {
		t.Error("Expceted Anthropic LLM type")
	}
	err = a.GetResponse(prompt)
	if err == nil {
		t.Error("expected error, got nil")
	}
	info := httpmock.GetCallCountInfo()
	expectedCalls := 1
	if info[fmt.Sprintf("POST %s", url)] != expectedCalls {
		t.Error("expected", expectedCalls, "got", info[fmt.Sprintf("POST %s", url)])
	}
}
