package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/gwenwindflower/tbd/shared"
	"github.com/schollz/progressbar/v3"
)

type Payload struct {
	Stop     interface{} `json:"stop"`
	Model    string      `json:"model"`
	Messages []Message   `json:"messages"`
	Temp     float64     `json:"temperature"`
	Tokens   int         `json:"max_tokens"`
	TopP     int         `json:"top_p"`
	Stream   bool        `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	SystemFingerprint interface{} `json:"system_fingerprint"`
	ID                string      `json:"id"`
	Object            string      `json:"object"`
	Model             string      `json:"model"`
	Choices           []struct {
		Logprobs interface{} `json:"logprobs"`
		Message  struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
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

// Groq API constants
const (
	maxRate     = 30
	interval    = time.Minute
	URL         = "https://api.groq.com/openai/v1/chat/completions"
	desc_prompt = `Generate a description for a column in a specific table in a data warehouse,
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
so again, no title, no formatting, just 1 to 3 sentences.`
	tests_prompt = `Generate a list of tests that can be run on a column in a specific table in a data warehouse,
the table is called %s and the column is called %s. The tests are YAML config, there are 2 to choose from.
They have the following structure, follow this structure exactly:
  - unique
  - not_null
Return only the tests that are applicable to the column, for example, a column that is a primary key should have 
both unique and not_null tests, while a column that is a foreign key should only have the not_null test. If a 
column is potentially optional, then it should have neither test. Return only the tests that are applicable to the column.
They will be nested under a 'tests' key in a YAML file, so no need to add a title or key, just the list of tests by themselves.
  For example, a good response for a 'product_type' column in an 'orders' table would be:
  - not_null

  A good response for an 'order_id' column in an 'orders' table would be:
  - unique
  - not_null

  A good response for a 'product_sku' column in an 'orders' table would be:
  - not_null
`
)

func GenerateColumnDescriptions(ts shared.SourceTables) {
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, maxRate)
	// We maek 2 calls so we divide the rate by 2
	limiter := time.NewTicker(interval / (maxRate / 2))
	defer limiter.Stop()

	bar := progressbar.NewOptions(len(ts.SourceTables),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("🤖📝"),
	)
	for i := range ts.SourceTables {
		for j := range ts.SourceTables[i].Columns {

			semaphore <- struct{}{}
			<-limiter.C

			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				table_name := ts.SourceTables[i].Name
				column_name := ts.SourceTables[i].Columns[j].Name
				desc_prompt := fmt.Sprintf(desc_prompt, table_name, column_name)
				tests_prompt := fmt.Sprintf(tests_prompt, table_name, column_name)
				desc_resp, err := GetGroqResponse(desc_prompt)
				if err != nil {
					log.Fatalf("Failed to get response from Groq for description: %v\n", err)
				}
				tests_resp, err := GetGroqResponse(tests_prompt)
				if err != nil {
					log.Fatalf("Failed to get response from Groq for tests: %v\n", err)
				}
				if len(desc_resp.Choices) > 0 {
					ts.SourceTables[i].Columns[j].Description = desc_resp.Choices[0].Message.Content
				}
				if len(tests_resp.Choices) > 0 {
					r := regexp.MustCompile(`unique|not_null`)
					matches := r.FindAllString(tests_resp.Choices[0].Message.Content, -1)
					matches = Deduplicate(matches)
					ts.SourceTables[i].Columns[j].Tests = matches
				}
			}(i, j)
		}
		bar.Add(1)
	}
	wg.Wait()
}

func GetGroqResponse(prompt string) (GroqResponse, error) {
	meta := Payload{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:  "llama3-70b-8192",
		Temp:   0.5,
		Tokens: 2048,
		TopP:   1,
		Stream: false,
		Stop:   nil,
	}
	payload, err := json.Marshal(meta)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v\n", err)
	}
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Unable to create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Cannot read response body: %v\n", err)
	}

	var resp GroqResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v\n", err)
	}
	return resp, nil
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
