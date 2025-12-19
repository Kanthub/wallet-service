package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GraphClient struct {
	client *http.Client
}

func NewGraphClient() *GraphClient {
	return &GraphClient{
		client: &http.Client{Timeout: 12 * time.Second},
	}
}

type graphReq struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type graphResp struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *GraphClient) Query(ctx context.Context, url string, query string, vars map[string]interface{}, out any) error {
	body, _ := json.Marshal(graphReq{Query: query, Variables: vars})
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("graph http status %d", resp.StatusCode)
	}

	var gr graphResp
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return err
	}
	if len(gr.Errors) > 0 {
		return fmt.Errorf("graph error: %s", gr.Errors[0].Message)
	}

	return json.Unmarshal(gr.Data, out)
}
