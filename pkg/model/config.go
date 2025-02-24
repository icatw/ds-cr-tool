package model

// DefaultConfig 默认的模型配置
var DefaultConfig = Config{
	Type:        "deepseek",
	Model:       "deepseek-ai/DeepSeek-R1",
	MaxTokens:   2000,
	Temperature: 0.7,
}

// NewConfig 创建新的配置实例
func NewConfig(apiKey string) *Config {
	cfg := DefaultConfig
	cfg.APIKey = apiKey
	return &cfg
}

// NewDsDefaultConfig 创建默认ds配置实例
func NewDsDefaultConfig(apiKey string) *Config {
	return &Config{
		Type:        "deepseek",
		Model:       "deepseek-ai/DeepSeek-R1",
		MaxTokens:   2000,
		Temperature: 0.7,
		APIKey:      apiKey,
	}
}
