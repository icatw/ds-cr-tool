package review

import "regexp"

// Rule 表示一条评审规则
type Rule struct {
	ID          string
	Name        string
	Description string
	Severity    string // "error", "warning", "info"
	Category    string // "security", "performance", "style", etc.
	Check       func(changes []FileChange) []Issue
}

// Issue 表示一个评审发现的问题
type Issue struct {
	Title       string
	FilePath    string
	Line        int
	Severity    string
	Description string
	Message     string
	Suggestion  string
	CodeSnippet string
}

// RuleEngine 规则引擎
type RuleEngine struct {
	rules []Rule
}

// NewRuleEngine 创建新的规则引擎
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: make([]Rule, 0),
	}
}

// AddRule 添加评审规则
func (e *RuleEngine) AddRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

// RunRules 执行所有评审规则
func (e *RuleEngine) RunRules(changes []FileChange) ([]Issue, error) {
	var issues []Issue

	for _, rule := range e.rules {
		ruleIssues := rule.Check(changes)
		issues = append(issues, ruleIssues...)
	}

	return issues, nil
}

// DefaultRules 返回默认的评审规则集
func DefaultRules() []Rule {
	return []Rule{
		{
			ID:          "SEC001",
			Name:        "敏感信息检查",
			Description: "检查代码中是否包含敏感信息（如密码、Token等）",
			Severity:    "error",
			Category:    "security",
			Check: func(changes []FileChange) []Issue {
				var issues []Issue
				for _, change := range changes {
					// 检查每一行代码
					for i, line := range change.Lines {
						// 检查敏感信息
						sensitivePatterns := map[string]string{
							"(?i)password.*=":   "避免在代码中硬编码密码",
							"(?i)api[_]?key.*=": "避免在代码中硬编码API密钥",
							"(?i)secret.*=":     "避免在代码中硬编码密钥",
							"(?i)token.*=":      "避免在代码中硬编码令牌",
						}

						for pattern, suggestion := range sensitivePatterns {
							matched, _ := regexp.MatchString(pattern, line)
							if matched {
								issues = append(issues, Issue{
									Title:       "发现敏感信息",
									FilePath:    change.FilePath,
									Line:        i + 1,
									Severity:    "error",
									Description: "代码中包含敏感信息",
									Message:     "检测到可能的敏感信息泄露",
									Suggestion:  suggestion,
									CodeSnippet: line,
								})
							}
						}
					}
				}
				return issues
			},
		},
		{
			ID:          "PERF001",
			Name:        "性能优化建议",
			Description: "检查可能影响性能的代码模式",
			Severity:    "warning",
			Category:    "performance",
			Check: func(changes []FileChange) []Issue {
				var issues []Issue
				for _, change := range changes {
					// 检查每一行代码
					for i, line := range change.Lines {
						// 检查可能的性能问题
						perfPatterns := map[string]string{
							"for.*range.*\\{$": "考虑使用带缓冲的channel或worker池来优化并发处理",
							"time\\.Sleep":     "避免使用time.Sleep，考虑使用定时器或超时控制",
							"append.*append":   "多次append可能导致频繁的内存分配，考虑预分配切片容量",
							"json\\.Marshal":   "在循环中进行JSON序列化可能影响性能，考虑复用编码器",
							"ioutil\\.ReadAll": "读取大文件时避免一次性读入内存，建议使用bufio.Scanner",
							"fmt\\.Sprintf":    "在性能敏感的场景中，考虑使用strings.Builder替代fmt.Sprintf",
						}

						for pattern, suggestion := range perfPatterns {
							matched, _ := regexp.MatchString(pattern, line)
							if matched {
								issues = append(issues, Issue{
									Title:       "发现性能优化机会",
									FilePath:    change.FilePath,
									Line:        i + 1,
									Severity:    "warning",
									Description: "发现可能影响性能的代码模式",
									Message:     "发现潜在的性能优化机会",
									Suggestion:  suggestion,
									CodeSnippet: line,
								})
							}
						}
					}
				}
				return issues
			},
		},
	}
}
