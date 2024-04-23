package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (o *OpenAI) GetResponse(prompt string) error {
	meta := Payload{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:     o.Model,
		Temp:      0.5,
		MaxTokens: 2048,
		TopP:      1,
		Stream:    false,
	}
	payload, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, o.Url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.ApiKey)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %v", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %v", err)
	}
	err = json.Unmarshal(body, &o.Response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}

func (o *Groq) GetResponse(prompt string) error {
	meta := Payload{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:     o.Model,
		Temp:      0.5,
		MaxTokens: 2048,
		TopP:      1,
		Stream:    false,
	}
	payload, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, o.Url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.ApiKey)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %v", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %v", err)
	}
	err = json.Unmarshal(body, &o.Response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}

func (a *Anthropic) GetResponse(prompt string) error {
	meta := Payload{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:     a.Model,
		Temp:      0.5,
		MaxTokens: 2048,
		TopP:      1,
		Stream:    false,
	}
	payload, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, a.Url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("x-api-key", a.ApiKey)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %v", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %v", err)
	}
	err = json.Unmarshal(body, &a.Response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}

func Deduplicate(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}
