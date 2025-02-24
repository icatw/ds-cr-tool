package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/icatw/ds-cr-tool/pkg/cache"
	"github.com/icatw/ds-cr-tool/pkg/cli"
	"github.com/icatw/ds-cr-tool/pkg/git"
	"github.com/icatw/ds-cr-tool/pkg/model"
	"github.com/icatw/ds-cr-tool/pkg/review"
)

func main() {
	// 解析命令行参数
	opts, err := cli.ParseFlags()
	if err != nil {
		log.Fatalf("解析参数失败: %v\n", err)
	}

	// 初始化Git客户端
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作目录失败: %v\n", err)
	}
	gitClient := git.NewGitClient(wd)
	if gitClient == nil {
		log.Fatalf("初始化Git客户端失败\n")
	}

	// 初始化代码分析器
	analyzer := review.NewAnalyzer(gitClient)

	// 获取代码改动
	var changes []review.FileChange
	if opts.CommitRange != "" {
		changes, err = analyzer.AnalyzeChanges(opts.CommitRange, "")
	} else {
		// 默认分析最新的改动
		changes, err = analyzer.AnalyzeChanges("HEAD~1", "HEAD")
	}
	if err != nil {
		log.Fatalf("分析代码改动失败: %v\n", err)
	}

	// 初始化缓存
	cacheDir := filepath.Join(os.Getenv("HOME"), ".cr", "cache")
	reviewCache, err := cache.NewReviewCache(cacheDir)
	if err != nil {
		log.Printf("初始化缓存失败: %v\n", err)
	}

	// 初始化AI模型客户端
	modelCfg := &model.Config{
		Type:   "deepseek",
		APIKey: os.Getenv("DEEPSEEK_API_KEY"),
		Model:  "deepseek-chat",
	}
	modelClient, err := model.NewModelClient(modelCfg)
	if err != nil {
		log.Fatalf("初始化AI模型客户端失败: %v\n", err)
	}

	// 处理每个改动文件
	for _, change := range changes {
		// 检查缓存
		if reviewCache != nil {
			if cached, err := reviewCache.Get(change.DiffContent); err == nil && cached != nil {
				fmt.Printf("使用缓存的评审结果 - %s\n", change.FilePath)
				fmt.Println(cached.ReviewResult)
				continue
			}
		}

		// 调用AI进行评审
		req := &model.ChatRequest{
			Messages: []model.Message{
				{
					Role:    "system",
					Content: "你是一个专业的代码评审助手，请对以下代码改动进行评审并提供建议。",
				},
				{
					Role:    "user",
					Content: fmt.Sprintf("文件: %s\n改动类型: %s\n\n%s", change.FilePath, change.ChangeType, change.DiffContent),
				},
			},
		}

		resp, err := modelClient.Chat(req)
		if err != nil {
			log.Printf("评审失败 - %s: %v\n", change.FilePath, err)
			continue
		}

		// 输出评审结果
		fmt.Printf("\n=== 文件: %s ===\n", change.FilePath)
		fmt.Println(resp.Choices[0].Message.Content)

		// 缓存评审结果
		if reviewCache != nil {
			expireAfter := 24 * time.Hour
			if err := reviewCache.Set(change.DiffContent, resp.Choices[0].Message.Content, &expireAfter); err != nil {
				log.Printf("缓存评审结果失败: %v\n", err)
			}
		}
	}
}