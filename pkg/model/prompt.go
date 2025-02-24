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

// DefaultReviewPrompt 返回默认的代码评审提示模板
func DefaultReviewPrompt() *ReviewPrompt {
    return &ReviewPrompt{
        BasePrompt: `你是一个专业的代码评审专家，请对以下代码变更进行全面的评审。

评审要求：
1. 使用中文进行评审
2. 保持客观专业的评审态度
3. 提供具体的改进建议
4. 重视代码质量和安全性
5. 关注性能优化机会

请重点关注：
1. 代码质量和最佳实践
   - 代码结构和组织
   - 命名规范
   - 注释完整性
   - 错误处理
2. 潜在的bug和安全问题
   - 边界条件处理
   - 并发安全
   - 资源管理
   - 输入验证
3. 性能优化机会
   - 算法效率
   - 资源使用
   - 缓存策略
4. 可维护性和可读性
   - 代码复杂度
   - 重复代码
   - 模块化
5. 测试覆盖情况
   - 单元测试
   - 边界测试
   - 错误场景测试

请按照以下格式提供评审意见：

## 总体评价
[简要总结代码变更的质量，包括主要优点和待改进点]

## 主要发现
[按重要性排序列出主要问题和建议]

## 详细分析
### 代码质量
[详细说明代码质量相关的问题和建议]

### 安全性
[详细说明安全相关的问题和建议]

### 性能
[详细说明性能相关的问题和建议]

### 可维护性
[详细说明可维护性相关的问题和建议]

### 测试建议
[详细说明测试相关的建议]

## 改进建议
[具体的改进建议，包括代码示例]`,
        FocusAreas: []string{
            "代码质量",
            "安全性",
            "性能",
            "可维护性",
            "测试覆盖",
        },
        OutputFormat: "markdown",
        LanguageBestPractices: map[string][]string{
            "go": {
                "使用 go fmt 格式化代码",
                "遵循 Go 官方代码规范",
                "适当使用接口",
                "正确处理错误",
                "使用 context 控制并发",
                "优化内存分配",
                "合理使用goroutine池",
            },
            "python": {
                "遵循PEP 8规范",
                "使用类型注解",
                "合理使用异常处理",
                "避免全局变量",
                "使用列表推导式优化性能",
            },
            "javascript": {
                "使用ESLint规范代码",
                "合理使用async/await",
                "避免内存泄漏",
                "使用现代ES特性",
                "注意类型安全",
            },
        },
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