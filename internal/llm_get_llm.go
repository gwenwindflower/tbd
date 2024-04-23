package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/gwenwindflower/tbd/shared"
)

type Llm interface {
	GetResponse(prompt string) error
	SetDescription(descPrompt string, ts shared.SourceTables, i, j int) error
	SetTests(descPrompt string, ts shared.SourceTables, i, j int) error
	GetRateLimiter() (chan struct{}, *time.Ticker)
}

type OpenAI struct {
	Type     string
	ApiKey   string
	Model    string
	Url      string
	Response struct {
		SystemFingerprint interface{} `json:"system_fingerprint"`
		Id                string      `json:"id"`
		Object            string      `json:"object"`
		Model             string      `json:"model"`
		Choices           []struct {
			Logprobs     interface{} `json:"logprobs"`
			Message      Message     `json:"message"`
			FinishReason string      `json:"finish_reason"`
			Index        int         `json:"index"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int     `json:"prompt_tokens"`
			PromptTime       float64 `json:"prompt_time"`
			CompletionTokens int     `json:"completion_tokens"`
			CompletionTime   float64 `json:"completion_time"`
			TotalTokens      int     `json:"total_tokens"`
			TotalTime        float64 `json:"total_time"`
		} `json:"usage"`
		Created int `json:"created"`
	}
}

type Groq struct {
	Type     string
	ApiKey   string
	Model    string
	Url      string
	Response struct {
		SystemFingerprint interface{} `json:"system_fingerprint"`
		Id                string      `json:"id"`
		Object            string      `json:"object"`
		Model             string      `json:"model"`
		Choices           []struct {
			Logprobs     interface{} `json:"logprobs"`
			Message      Message     `json:"message"`
			FinishReason string      `json:"finish_reason"`
			Index        int         `json:"index"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int     `json:"prompt_tokens"`
			PromptTime       float64 `json:"prompt_time"`
			CompletionTokens int     `json:"completion_tokens"`
			CompletionTime   float64 `json:"completion_time"`
			TotalTokens      int     `json:"total_tokens"`
			TotalTime        float64 `json:"total_time"`
		} `json:"usage"`
		Created int `json:"created"`
	}
}

type Anthropic struct {
	Type     string
	ApiKey   string
	Model    string
	Url      string
	Response struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Id           string   `json:"id"`
		Type         string   `json:"type"`
		Role         string   `json:"role"`
		StopReason   string   `json:"stop_reason"`
		StopSequence []string `json:"stop_sequence"`
		Usage        struct {
			InputTokens int `json:"input_tokens"`
			OutputToken int `json:"output_tokens"`
		} `json:"usage"`
	}
}

type Payload struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Temp      float64   `json:"temperature"`
	MaxTokens int       `json:"max_tokens"`
	TopP      int       `json:"top_p"`
	Stream    bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func GetLlm(fr FormResponse) (Llm, error) {
	switch fr.Llm {
	case "groq":
		g := &Groq{
			Type:   "groq",
			ApiKey: os.Getenv(fr.LlmKeyEnvVar),
			Model:  "llama3-70b-8192",
			Url:    "https://api.groq.com/openai/v1/chat/completions",
		}
		return g, nil
	case "openai":
		o := &OpenAI{
			Type:   "openai",
			ApiKey: os.Getenv(fr.LlmKeyEnvVar),
			Model:  "gpt-4-turbo",
			Url:    "https://api.openai.com/v1/chat/completions",
		}
		return o, nil
	case "anthropic":
		a := &Anthropic{
			Type:   "anthropic",
			ApiKey: os.Getenv(fr.LlmKeyEnvVar),
			Model:  "claude-3-opus-20240229",
			Url:    "https://api.anthropic.com/v1/messages",
		}
		return a, nil
	default:
		return nil, fmt.Errorf("invalid LLM choice: %v", fr.Llm)
	}
}
