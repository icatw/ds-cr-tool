package review

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// ReportFormat 定义支持的报告格式
type ReportFormat string

const (
	MarkdownFormat ReportFormat = "markdown"
	HTMLFormat    ReportFormat = "html"
	PDFFormat     ReportFormat = "pdf"
)

// Reporter 定义报告生成器接口
type Reporter interface {
	Generate(issues []Issue, format ReportFormat) ([]byte, error)
}

// DefaultReporter 默认报告生成器实现
type DefaultReporter struct {
	ProjectName string
	CommitID    string
}

// NewReporter 创建新的报告生成器
func NewReporter(projectName, commitID string) Reporter {
	return &DefaultReporter{
		ProjectName: projectName,
		CommitID:    commitID,
	}
}

// Generate 生成评审报告
func (r *DefaultReporter) Generate(issues []Issue, format ReportFormat) ([]byte, error) {
	switch format {
	case MarkdownFormat:
		return r.generateMarkdown(issues)
	case HTMLFormat:
		return nil, fmt.Errorf("HTML format not implemented yet")
	case PDFFormat:
		return nil, fmt.Errorf("PDF format not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// generateMarkdown 生成Markdown格式的报告
func (r *DefaultReporter) generateMarkdown(issues []Issue) ([]byte, error) {
	var buf bytes.Buffer

	// 写入报告头部
	buf.WriteString(fmt.Sprintf("# 代码评审报告\n\n"))
	buf.WriteString(fmt.Sprintf("## 项目信息\n\n"))
	buf.WriteString(fmt.Sprintf("- 项目名称：%s\n", r.ProjectName))
	buf.WriteString(fmt.Sprintf("- 提交ID：%s\n", r.CommitID))
	buf.WriteString(fmt.Sprintf("- 评审时间：%s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// 按严重程度分类统计
	severityCount := make(map[string]int)
	for _, issue := range issues {
		severityCount[string(issue.Severity)]++
	}

	// 写入统计信息
	buf.WriteString("## 评审结果统计\n\n")

	// 添加代码统计信息
	buf.WriteString("### 代码变更统计\n\n")
	buf.WriteString("| 指标 | 数值 |\n")
	buf.WriteString("|------|---------|\n")
	buf.WriteString(fmt.Sprintf("| 评审文件数 | %d |\n", len(getUniqueFiles(issues))))
	buf.WriteString(fmt.Sprintf("| 问题总数 | %d |\n", len(issues)))

	// 写入严重程度统计
	buf.WriteString("\n### 问题严重程度分布\n\n")
	buf.WriteString("| 严重程度 | 数量 |\n")
	buf.WriteString("|---------|---------|\n")
	for severity, count := range severityCount {
	    buf.WriteString(fmt.Sprintf("| %s | %d |\n", severity, count))
	}
	buf.WriteString("\n")

	// 写入优化建议总结
	buf.WriteString("## 整体优化建议\n\n")
	suggestions := summarizeSuggestions(issues)
	for _, suggestion := range suggestions {
	    buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
	}
	buf.WriteString("\n")

	// 写入详细问题列表
	buf.WriteString("## 详细问题列表\n\n")
	for i, issue := range issues {
		buf.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, issue.Title))
		buf.WriteString(fmt.Sprintf("- 文件：%s\n", issue.FilePath))
		buf.WriteString(fmt.Sprintf("- 位置：第%d行\n", issue.Line))
		buf.WriteString(fmt.Sprintf("- 严重程度：%s\n", issue.Severity))
		buf.WriteString(fmt.Sprintf("- 描述：%s\n", issue.Description))
		if issue.Suggestion != "" {
			buf.WriteString(fmt.Sprintf("- 建议：%s\n", issue.Suggestion))
		}
		buf.WriteString("\n")

		// 添加代码片段（如果有）
		if issue.CodeSnippet != "" {
			// 获取代码片段的上下文
			lines := strings.Split(issue.CodeSnippet, "\n")
			contextStart := max(0, issue.Line-3)
			contextEnd := min(len(lines), issue.Line+3)

			buf.WriteString("```go\n")
			// 添加行号
			for i := contextStart; i < contextEnd; i++ {
				linePrefix := "  "
				if i == issue.Line-1 { // 高亮问题行
					linePrefix = ">"
				}
				buf.WriteString(fmt.Sprintf("%s %4d | %s\n", linePrefix, i+1, lines[i]))
			}
			buf.WriteString("```\n\n")
		}
	}

	return buf.Bytes(), nil
}

// 辅助函数
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 辅助函数：获取唯一文件列表
func getUniqueFiles(issues []Issue) []string {
    filesMap := make(map[string]bool)
    for _, issue := range issues {
        filesMap[issue.FilePath] = true
    }
    
    files := make([]string, 0, len(filesMap))
    for file := range filesMap {
        files = append(files, file)
    }
    return files
}

// 辅助函数：汇总优化建议
func summarizeSuggestions(issues []Issue) []string {
    suggestionMap := make(map[string]int)
    for _, issue := range issues {
        if issue.Suggestion != "" {
            suggestionMap[issue.Suggestion]++
        }
    }
    
    suggestions := make([]string, 0, len(suggestionMap))
    for suggestion, count := range suggestionMap {
        if count > 1 {
            suggestions = append(suggestions, fmt.Sprintf("%s (出现%d次)", suggestion, count))
        } else {
            suggestions = append(suggestions, suggestion)
        }
    }
    return suggestions
}