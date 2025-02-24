package review

import (
	"fmt"
)

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
	RuleID      string
	FilePath    string
	LineNumber  int
	Severity    string
	Message     string
	Suggestion  string
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
				// TODO: 实现具体的检查逻辑
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
				// TODO: 实现具体的检查逻辑
				return issues
			},
		},
	}
}