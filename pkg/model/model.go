package model

// ModelClient 定义通用的AI模型客户端接口
type ModelClient interface {
	// Chat 发送聊天请求并获取响应
	Chat(req *ChatRequest) (*ChatResponse, error)
}

// Config 定义模型配置
type Config struct {
	// 模型类型（如 "deepseek", "openai" 等）
	Type string `json:"type"`
	// API密钥
	APIKey string `json:"api_key"`
	// 模型名称
	Model string `json:"model"`
	// 其他通用配置参数
	MaxTokens int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// NewModelClient 根据配置创建对应的模型客户端
func NewModelClient(cfg *Config) (ModelClient, error) {
	switch cfg.Type {
	case "deepseek":
		return NewDeepSeekClient(cfg.APIKey), nil
	default:
		return nil, fmt.Errorf("unsupported model type: %s", cfg.Type)
	}
}