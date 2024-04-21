package main

import (
	"fmt"
	"regexp"

	"github.com/gwenwindflower/tbd/shared"
)

func (o *OpenAI) SetTests(testsPrompt string, ts shared.SourceTables, i, j int) error {
	err := o.GetResponse(testsPrompt)
	if err != nil {
		return fmt.Errorf("failed to get response from OpenAI for tests: %v", err)
	}
	if len(o.Response.Choices) > 0 {
		r := regexp.MustCompile(`unique|not_null`)
		matches := r.FindAllString(o.Response.Choices[0].Message.Content, -1)
		matches = Deduplicate(matches)
		ts.SourceTables[i].Columns[j].Tests = matches
	}
	return nil
}

func (g *Groq) SetTests(testsPrompt string, ts shared.SourceTables, i, j int) error {
	err := g.GetResponse(testsPrompt)
	if err != nil {
		return fmt.Errorf("failed to get response from Groq for tests: %v", err)
	}
	if len(g.Response.Choices) > 0 {
		r := regexp.MustCompile(`unique|not_null`)
		matches := r.FindAllString(g.Response.Choices[0].Message.Content, -1)
		matches = Deduplicate(matches)
		ts.SourceTables[i].Columns[j].Tests = matches
	}
	return nil
}

func (a *Anthropic) SetTests(testsPrompt string, ts shared.SourceTables, i, j int) error {
	err := a.GetResponse(testsPrompt)
	if err != nil {
		return fmt.Errorf("failed to get response from Anthropic for tests: %v", err)
	}
	if len(a.Response.Content) == 0 {
		return fmt.Errorf("no response content, likely bad request")
	}
	resp := a.Response.Content[0]
	if len(resp.Text) > 0 && resp.Type == "text" {
		r := regexp.MustCompile(`unique|not_null`)
		matches := r.FindAllString(resp.Text, -1)
		matches = Deduplicate(matches)
		ts.SourceTables[i].Columns[j].Tests = matches
	}
	return nil
}
