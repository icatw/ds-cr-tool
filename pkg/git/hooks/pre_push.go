package hooks

import (
	"fmt"
	"os"
	"strings"

	"github.com/icatw/ds-cr-tool/pkg/git"
	"github.com/icatw/ds-cr-tool/pkg/model"
	"github.com/icatw/ds-cr-tool/pkg/review"
)

// PrePushHook 处理pre-push钩子的逻辑
type PrePushHook struct {
	Options map[string]string
}

// NewPrePushHook 创建新的pre-push钩子处理器
func NewPrePushHook(options map[string]string) *PrePushHook {
	return &PrePushHook{
		Options: options,
	}
}

// Execute 执行pre-push钩子逻辑
func (h *PrePushHook) Execute() error {
	// 获取标准输入中的引用信息
	refInfo, err := h.readRefInfo()
	if err != nil {
		return fmt.Errorf("failed to read ref info: %v", err)
	}

	// 解析引用信息
	refs := h.parseRefInfo(refInfo)
	if len(refs) == 0 {
		return fmt.Errorf("no refs to push")
	}

	// 对每个要推送的引用进行代码评审
	for _, ref := range refs {
		if err := h.reviewRef(ref); err != nil {
			return fmt.Errorf("review failed for ref %s: %v", ref.Name, err)
		}
	}

	return nil
}

// RefInfo 存储引用信息
type RefInfo struct {
	Name     string
	OldHash  string
	NewHash  string
	Remote   string
}

// readRefInfo 从标准输入读取引用信息
func (h *PrePushHook) readRefInfo() (string, error) {
	info, err := os.ReadFile(os.Stdin.Name())
	if err != nil {
		return "", err
	}
	return string(info), nil
}

// parseRefInfo 解析引用信息
func (h *PrePushHook) parseRefInfo(info string) []RefInfo {
	var refs []RefInfo
	lines := strings.Split(info, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 4 {
			continue
		}

		refs = append(refs, RefInfo{
			Name:    parts[2],
			OldHash: parts[0],
			NewHash: parts[1],
			Remote:  parts[3],
		})
	}

	return refs
}

// reviewRef 对指定引用进行代码评审
func (h *PrePushHook) reviewRef(ref RefInfo) error {
	// 如果是删除分支操作，则跳过评审
	if ref.NewHash == "0000000000000000000000000000000000000000" {
		return nil
	}

	// 创建Git客户端
	gitClient := git.NewGitClient(h.Options["repo_path"])

	// 创建代码分析器
	analyzer := review.NewAnalyzer(gitClient)

	// 获取改动的文件列表和差异内容
	changes, err := analyzer.AnalyzeChanges(ref.OldHash, ref.NewHash)
	if err != nil {
		return fmt.Errorf("分析代码改动失败: %v", err)
	}

	// 如果没有改动，直接返回
	if len(changes) == 0 {
		return nil
	}

	// 创建模型客户端
	modelCfg := &model.Config{
		Type: "deepseek",
		APIKey: h.Options["api_key"],
		Model: "deepseek-coder",
		MaxTokens: 2048,
		Temperature: 0.7,
	}

	modelClient, err := model.NewModelClient(modelCfg)
	if err != nil {
		return fmt.Errorf("创建模型客户端失败: %v", err)
	}

	// 创建评审报告生成器
	reporter := review.NewReporter("ds-cr-tool", ref.NewHash)

	// 创建缓存管理器
	cacheManager, err := cache.NewReviewCache(h.Options["cache_dir"])
	if err != nil {
		return fmt.Errorf("创建缓存管理器失败: %v", err)
	}

	// 分析代码问题
	issues := make([]review.Issue, 0)
	for _, change := range changes {
		// 检查缓存
		cacheItem, err := cacheManager.Get(change.DiffContent)
		if err != nil {
			return fmt.Errorf("读取缓存失败: %v", err)
		}

		var reviewResult string
		if cacheItem != nil {
			// 使用缓存的评审结果
			reviewResult = cacheItem.ReviewResult
		} else {
			// 使用模型分析代码问题
			response, err := modelClient.Analyze(change.DiffContent)
			if err != nil {
				return fmt.Errorf("分析代码问题失败: %v", err)
			}

			// 将分析结果转换为JSON字符串
			reviewResult, err = json.Marshal(response)
			if err != nil {
				return fmt.Errorf("序列化评审结果失败: %v", err)
			}

			// 缓存评审结果
			expireAfter := 24 * time.Hour
			if err := cacheManager.Set(change.DiffContent, string(reviewResult), &expireAfter); err != nil {
				return fmt.Errorf("缓存评审结果失败: %v", err)
			}
		}

		// 解析评审结果
		var response model.AnalyzeResponse
		if err := json.Unmarshal([]byte(reviewResult), &response); err != nil {
			return fmt.Errorf("解析评审结果失败: %v", err)
		}

		// 将模型分析结果转换为Issue
		for _, item := range response.Items {
			issues = append(issues, review.Issue{
				Title:      item.Title,
				FilePath:   change.FilePath,
				Severity:   review.Severity(item.Severity),
				Message:    item.Message,
				Suggestion: item.Suggestion,
			})
		}
	}

	// 生成评审报告
	reportBytes, err := reporter.Generate(issues, review.MarkdownFormat)
	if err != nil {
		return fmt.Errorf("生成评审报告失败: %v", err)
	}

	// 检查是否有严重问题
	hasCriticalIssues := false
	for _, issue := range issues {
		if issue.Severity == review.SeverityCritical {
			hasCriticalIssues = true
			break
		}
	}

	// 如果有严重问题，阻止推送
	if hasCriticalIssues {
		return fmt.Errorf("代码评审未通过，存在严重问题:\n%s", string(reportBytes))
	}

	return nil
}