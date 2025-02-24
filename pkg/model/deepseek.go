package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	MaxTokens int      `json:"max_tokens"`
	Stop     []string  `json:"stop,omitempty"`
	Temperature float64 `json:"temperature"`
	TopP     float64   `json:"top_p"`
	TopK     int       `json:"top_k"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	N        int       `json:"n"`
	ResponseFormat map[string]string `json:"response_format"`
	Tools    []Tool    `json:"tools,omitempty"`
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
	ID        string    `json:"id"`
	Choices   []Choice  `json:"choices"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Usage     Usage     `json:"usage"`
	Created   int64     `json:"created"`
	Model     string    `json:"model"`
	Object    string    `json:"object"`
}

// Choice 定义响应选项的结构
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
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

// Usage 定义token使用情况的结构
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
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

	body, err := ioutil.ReadAll(resp.Body)
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