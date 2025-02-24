package model

// ReviewPrompt 定义代码评审的提示模板
type ReviewPrompt struct {
    // 基础提示信息
    BasePrompt string
    // 评审重点
    FocusAreas []string
    // 输出格式
    OutputFormat string
    // 语言相关的最佳实践
    LanguageBestPractices map[string][]string
}

// GeneratePrompt 根据代码差异生成完整的评审提示
func (p *ReviewPrompt) GeneratePrompt(diff string) []Message {
	return []Message{
		{
			Role:    "system",
			Content: p.BasePrompt,
		},
		{
			Role:    "user",
			Content: "请评审以下代码变更：\n\n" + diff,
		},
	}
}