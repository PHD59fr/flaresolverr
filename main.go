package flaresolverr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Init(i ClientInput) (*Client, error) {
	client := &Client{
		BaseUrl:       i.BaseUrl,
		TimeOut:       i.TimeOut,
		CustomHeaders: i.CustomHeaders,
	}

	resp, err := http.Get(client.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to FlareSolverr: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Errorf("F")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FlareSolverr returned HTTP status: %d", resp.StatusCode)
	}

	var data struct {
		Msg       string `json:"msg"`
		Version   string `json:"version"`
		UserAgent string `json:"userAgent"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding FlareSolverr response: %w", err)
	}

	if data.Msg != "FlareSolverr is ready!" {
		return nil, fmt.Errorf("FlareSolverr is not initialized: %s", data.Msg)
	}

	client.Version = data.Version
	client.UserAgent = data.UserAgent
	client.BaseUrl = i.BaseUrl + "/v1"
	if i.UserAgent != "" {
		client.UserAgent = i.UserAgent
	}
	client.TimeOut = 60000
	if i.TimeOut != 0 {
		client.TimeOut = i.TimeOut
	}
	return client, nil
}

func (c *Client) Get(targetUrl string) (*Response, error) {
	return c.sendRequest("request.get", targetUrl)
}

func (c *Client) Post(targetUrl string, data map[string]interface{}) (*Response, error) {
	return c.sendRequestWithData("request.post", targetUrl, data)
}

func (c *Client) Put(targetUrl string, data map[string]interface{}) (*Response, error) {
	return c.sendRequestWithData("request.put", targetUrl, data)
}

func (c *Client) Delete(targetUrl string) (*Response, error) {
	return c.sendRequest("request.delete", targetUrl)
}

func (c *Client) sendRequest(cmd, targetUrl string) (*Response, error) {
	payload := map[string]interface{}{
		"cmd":        cmd,
		"url":        targetUrl,
		"maxTimeout": c.TimeOut,
	}
	return c.sendRequestPayload(payload)
}

func (c *Client) sendRequestWithData(cmd, targetUrl string, data map[string]interface{}) (*Response, error) {
	payload := map[string]interface{}{
		"cmd":        cmd,
		"url":        targetUrl,
		"maxTimeout": c.TimeOut,
		"data":       data,
	}
	return c.sendRequestPayload(payload)
}

func (c *Client) sendRequestPayload(payload map[string]interface{}) (*Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	for key, value := range c.CustomHeaders {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FlareSolverr returned HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	bodyResponse := &Response{}
	if err := json.Unmarshal(body, bodyResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return bodyResponse, nil
}
