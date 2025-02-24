package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DeepSeekAPIURL = "https://api.siliconflow.cn/v1/chat/completions"
)

// DeepSeekClient 实现DeepSeek API的客户端
type DeepSeekClient struct {
	config *Config
	client *http.Client
}

// NewDeepSeekClient 创建新的DeepSeek客户端实例
func NewDeepSeekClient(cfg *Config) *DeepSeekClient {
	return &DeepSeekClient{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Chat 发送聊天请求并获取响应
func (c *DeepSeekClient) Chat(req *ChatRequest) (*ChatResponse, error) {
	// 使用配置参数
	if req.Model == "" {
		req.Model = c.config.Model
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = c.config.MaxTokens
	}
	if req.Temperature == 0 {
		req.Temperature = c.config.Temperature
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %v", err)
	}

	// 添加重试机制
	var resp *http.Response
	var lastErr error
	for retries := 0; retries < 3; retries++ {
		httpReq, err := http.NewRequest("POST", DeepSeekAPIURL, bytes.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("create request failed: %v", err)
		}

		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err = c.client.Do(httpReq)
		if err == nil {
			break
		}
		lastErr = err
		time.Sleep(time.Duration(retries+1) * time.Second)
	}

	if resp == nil {
		return nil, fmt.Errorf("all retries failed: %v", lastErr)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &chatResp, nil
}

// ChatRequest 定义聊天请求的参数结构
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	MaxTokens   int       `json:"max_tokens"`
	Stop        []string  `json:"stop,omitempty"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	//TopK             int               `json:"top_k"`
	FrequencyPenalty float64           `json:"frequency_penalty"`
	N                int               `json:"n"`
	ResponseFormat   map[string]string `json:"response_format"`
	Tools            []Tool            `json:"tools,omitempty"`
}

// Message 定义聊天消息的结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool 定义工具的结构
type Tool struct {
	Type     string         `json:"type"`
	Function FunctionConfig `json:"function"`
}

// FunctionConfig 定义函数配置的结构
type FunctionConfig struct {
	Description string                 `json:"description"`
	Name        string                 `json:"name"`
	Parameters  map[string]interface{} `json:"parameters"`
	Strict      bool                   `json:"strict"`
}

// ChatResponse 定义API响应的结构
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 定义响应选项的结构
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 定义token使用情况的结构
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ToolCall 定义工具调用的结构
type ToolCall struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Function FunctionResult `json:"function"`
}

// FunctionResult 定义函数调用结果的结构
type FunctionResult struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
