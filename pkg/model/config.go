package model

// DefaultModelConfig 默认的全局模型配置
var DefaultModelConfig = ModelConfig{
	DefaultModel: "deepseek",
	Models: map[string]*Config{
		"deepseek": {
			Type:        "deepseek",
			Model:       "deepseek-ai/DeepSeek-R1",
			MaxTokens:   2000,
			Temperature: 0.7,
		},
		"openai": {
			Type:        "openai",
			Model:       "gpt-3.5-turbo",
			MaxTokens:   2000,
			Temperature: 0.7,
		},
		"chatglm": {
			Type:        "chatglm",
			Model:       "glm-4",
			MaxTokens:   2000,
			Temperature: 0.7,
		},
		"qwen": {
			Type:        "qwen",
			Model:       "qwen-coder-plus",
			MaxTokens:   2000,
			Temperature: 0.7,
		},
	},
}

// NewModelConfigWithKeys 创建带有API密钥的模型配置
func NewModelConfigWithKeys(deepseekKey, openaiKey, chatglmKey, qwenKey string) *ModelConfig {
	cfg := DefaultModelConfig
	if deepseekKey != "" {
		cfg.Models["deepseek"].APIKey = deepseekKey
	}
	if openaiKey != "" {
		cfg.Models["openai"].APIKey = openaiKey
	}
	if chatglmKey != "" {
		cfg.Models["chatglm"].APIKey = chatglmKey
	}
	if qwenKey != "" {
		cfg.Models["qwen"].APIKey = qwenKey
	}
	return &cfg
}

// NewConfig 创建新的配置实例（保留用于兼容性）
func NewConfig(apiKey string) *Config {
	return &Config{
		Type:        "deepseek",
		Model:       "deepseek-ai/DeepSeek-R1",
		MaxTokens:   2000,
		Temperature: 0.7,
		APIKey:      apiKey,
	}
}
