package internal

import (
	"fmt"

	"github.com/gwenwindflower/tbd/shared"
)

func (o *OpenAI) SetDescription(descPrompt string, ts shared.SourceTables, i, j int) error {
	err := o.GetResponse(descPrompt)
	if err != nil {
		return fmt.Errorf("failed to get response from OpenAI for description: %v", err)
	}
	if len(o.Response.Choices) > 0 {
		ts.SourceTables[i].Columns[j].Description = o.Response.Choices[0].Message.Content
	}

	return nil
}

func (g *Groq) SetDescription(descPrompt string, ts shared.SourceTables, i, j int) error {
	err := g.GetResponse(descPrompt)
	if err != nil {
		return fmt.Errorf("failed to get response from Groq for description: %v", err)
	}
	if len(g.Response.Choices) > 0 {
		ts.SourceTables[i].Columns[j].Description = g.Response.Choices[0].Message.Content
	}
	return nil
}

func (a *Anthropic) SetDescription(descPrompt string, ts shared.SourceTables, i, j int) error {
	err := a.GetResponse(descPrompt)
	if err != nil {
		return fmt.Errorf("failed to get ok response from Anthropic for description: %v", err)
	}
	if len(a.Response.Content) == 0 {
		return fmt.Errorf("no response content, likely bad request")
	}
	respContent := a.Response.Content[0]
	ts.SourceTables[i].Columns[j].Description = respContent.Text
	return nil
}
