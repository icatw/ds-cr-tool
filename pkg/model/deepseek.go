package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	DeepSeekAPIURL = "https://api.siliconflow.cn/v1/chat/completions"
)

// DeepSeekClient 实现DeepSeek API的客户端
type DeepSeekClient struct {
	apiKey string
}

// NewDeepSeekClient 创建新的DeepSeek客户端实例
func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey: apiKey,
	}
}

// ChatRequest 定义聊天请求的参数结构
type ChatRequest struct {
	Model            string            `json:"model"`
	Messages         []Message         `json:"messages"`
	Stream           bool              `json:"stream"`
	MaxTokens        int               `json:"max_tokens"`
	Stop             []string          `json:"stop,omitempty"`
	Temperature      float64           `json:"temperature"`
	TopP             float64           `json:"top_p"`
	TopK             int               `json:"top_k"`
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

// ReviewCode 发送代码评审请求并处理响应
func (c *DeepSeekClient) ReviewCode(prompt []Message) (*ChatResponse, error) {
	req := ChatRequest{
		Model:       "deepseek-ai/DeepSeek-R1",
		Messages:    prompt,
		MaxTokens:   2000,
		Temperature: 0.7,
		TopP:        0.95,
		Stream:      false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequest("POST", DeepSeekAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败: %s", string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &chatResp, nil
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

// Chat 发送聊天请求并获取响应
func (c *DeepSeekClient) Chat(req *ChatRequest) (*ChatResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %v", err)
	}

	httpReq, err := http.NewRequest("POST", DeepSeekAPIURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
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
