package model

// ReviewPrompt 定义代码评审的提示模板
type ReviewPrompt struct {
	// 基础提示信息
	BasePrompt string
	// 评审重点
	FocusAreas []string
	// 输出格式
	OutputFormat string
}

// DefaultReviewPrompt 返回默认的代码评审提示模板
func DefaultReviewPrompt() *ReviewPrompt {
	return &ReviewPrompt{
		BasePrompt: `你是一个专业的代码评审专家，请对以下代码变更进行全面的评审。
请重点关注：
1. 代码质量和最佳实践
2. 潜在的bug和安全问题
3. 性能优化机会
4. 可维护性和可读性
5. 测试覆盖情况

请按照以下格式提供评审意见：

## 总体评价
[对代码变更的整体评价]

## 主要发现
[列出主要问题和建议]

## 详细分析
[按文件分类的详细评审意见]

## 改进建议
[具体的改进建议和最佳实践推荐]`,
		FocusAreas: []string{
			"代码质量",
			"安全性",
			"性能",
			"可维护性",
			"测试覆盖",
		},
		OutputFormat: "markdown",
	}
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