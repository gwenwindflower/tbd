package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Payload struct {
	Messages []Message   `json:"messages"`
	Model    string      `json:"model"`
	Temp     float64     `json:"temperature"`
	Tokens   int         `json:"max_tokens"`
	TopP     int         `json:"top_p"`
	Stream   bool        `json:"stream"`
	Stop     interface{} `json:"stop"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

// Groq API constants
const (
	maxRate  = 30
	interval = time.Minute
	URL      = "https://api.groq.com/openai/v1/chat/completions"
)

func GenerateColumnDescriptions(tables SourceTables) {
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, maxRate)
	limiter := time.NewTicker(interval / maxRate)
	defer limiter.Stop()

	for i := range tables.SourceTables {
		for j := range tables.SourceTables[i].Columns {

			semaphore <- struct{}{}
			<-limiter.C

			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				table_name := tables.SourceTables[i].Name
				column_name := tables.SourceTables[i].Columns[j].Name
				prompt := fmt.Sprintf(`Generate a description for a column in a specific table in a data warehouse,
  the table is called %s and the column is called %s. The description should be concise, 1 to 3 sentences,
  and inform both business users and technical data analyts about the purpose and contents of the column.
  Avoid using the column name in the description, as it is redundant — put another way do not use tautological
  descriptions, for example on an 'order_id' column saying "This is the id of an order". Don't do that. A good
  example for an 'order_id' column would be something like "This is the primary key of the orders table,
  each distinct order has a unique 'order_id'". Another good example for an orders table would be describing
  'product_type' as "The category of product, the bucket that a product falls into, for example 'electronics' or 'clothing'".
  Avoid making assumptions about the data, as you don't have access to it. Don't make assertions about data that you 
  haven't seen, just use business context, the table name, and the column to generate the description. The description.
  There is no need to add a title just the sentences that compose the description, it's being put onto a field in a YAML file, 
so again, no title, no formatting, just 1 to 3 sentences.`,
					table_name, column_name)

				meta := Payload{
					Messages: []Message{
						{
							Role:    "user",
							Content: prompt,
						},
					},
					Model:  "gemma-7b-it",
					Temp:   0.5,
					Tokens: 1024,
					TopP:   1,
					Stream: false,
					Stop:   nil,
				}
				payload, err := json.Marshal(meta)
				if err != nil {
					log.Fatalf("Failed to marshal JSON: %v", err)
				}
				req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(payload))
				if err != nil {
					log.Fatalf("Unable to create request: %v", err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
				client := http.Client{}
				response, err := client.Do(req)
				if err != nil {
					log.Fatalf("Request failed: %v", err)
				}
				defer response.Body.Close()

				body, err := io.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Cannot read response body: %v", err)
				}

				var resp GroqResponse
				err = json.Unmarshal(body, &resp)
				if err != nil {
					log.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				if len(resp.Choices) > 0 {
					tables.SourceTables[i].Columns[j].Description = resp.Choices[0].Message.Content
				}
			}(i, j)
		}
	}
}
