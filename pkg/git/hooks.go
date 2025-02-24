package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HookType 定义Git钩子类型
type HookType string

const (
	PreCommitHook HookType = "pre-commit"
	PrePushHook   HookType = "pre-push"
)

// HookManager Git钩子管理器
type HookManager struct {
	repoPath string
}

// NewHookManager 创建新的钩子管理器
func NewHookManager(repoPath string) *HookManager {
	return &HookManager{
		repoPath: repoPath,
	}
}

// InstallHook 安装Git钩子
func (m *HookManager) InstallHook(hookType HookType) error {
	hookPath := filepath.Join(m.repoPath, ".git", "hooks", string(hookType))

	// 检查钩子文件是否已存在
	if _, err := os.Stat(hookPath); err == nil {
		return fmt.Errorf("hook already exists: %s", hookPath)
	}

	// 创建钩子脚本内容
	content := m.generateHookScript(hookType)

	// 写入钩子文件
	if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("failed to write hook file: %v", err)
	}

	return nil
}

// RemoveHook 移除Git钩子
func (m *HookManager) RemoveHook(hookType HookType) error {
	hookPath := filepath.Join(m.repoPath, ".git", "hooks", string(hookType))

	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("failed to remove hook: %v", err)
	}

	return nil
}

// generateHookScript 生成钩子脚本内容
func (m *HookManager) generateHookScript(hookType HookType) string {
	var script strings.Builder

	// 添加脚本头部
	script.WriteString("#!/bin/sh\n\n")

	// 根据钩子类型生成不同的脚本内容
	switch hookType {
	case PreCommitHook:
		script.WriteString("# 运行代码评审工具\n")
		script.WriteString("ds-cr-tool review pre-commit || exit 1\n")
	case PrePushHook:
		script.WriteString("# 运行代码评审工具\n")
		script.WriteString("ds-cr-tool review pre-push || exit 1\n")
	}

	return script.String()
}