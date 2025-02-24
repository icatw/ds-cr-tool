# GPT-CodeReview 安装指南

## 安装要求

- Go 1.16 或更高版本
- Git

## 安装步骤

1. 从GitHub安装:

```bash
go install github.com/icatw/ds-cr-tool/cmd/cr@latest
```

2. 配置DeepSeek API密钥:

你可以通过以下两种方式之一配置API密钥：

### 方式1：环境变量

```bash
export DEEPSEEK_API_KEY=your_api_key
```

### 方式2：配置文件

创建配置文件 `~/.cr/config.json`：

```bash
mkdir -p ~/.cr
echo '{"api_key": "your_api_key"}' > ~/.cr/config.json
```

## 在项目中使用

### 1. 安装Git钩子

在你的项目根目录下执行：

```bash
# 创建pre-commit钩子
mkdir -p .git/hooks
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh

set -e

# 日志函数
log() {
    echo "[ds-cr-tool] $1"
}

# 获取暂存区的文件列表
log "正在获取暂存区文件列表..."
files=$(git diff --cached --name-only --diff-filter=ACM)

# 如果没有文件要提交，直接退出
if [ -z "$files" ]; then
    log "没有需要评审的文件，跳过检查"
    exit 0
fi

# 运行代码评审工具
log "开始运行代码评审..."
review_result=$(cr review --files "$files" 2>&1) || {
    log "代码评审失败，请修复以下问题："
    echo "$review_result" | sed 's/^/    /'
    exit 1
}

# 评审通过
log "代码评审通过 ✓"
exit 0
EOF

# 添加执行权限
chmod +x .git/hooks/pre-commit
```

### 2. 基本使用

1. 直接评审最新改动（无需参数）：

```bash
cr diff
```

2. 评审指定的文件：

```bash
cr diff --files=main.go,utils.go
```

3. 评审指定范围的提交：

```bash
cr review --commit-range=HEAD~3..HEAD
```

4. 输出HTML格式报告：

```bash
cr diff --format=html --output=review.html
```

### 命令参数

- `--files`: 指定要评审的文件列表，多个文件用逗号分隔
- `--commit-range`: 指定要评审的提交范围
- `--format`: 输出格式（markdown/html/pdf），默认为markdown
- `--output`: 输出文件路径，默认输出到标准输出
- `--verbose`: 显示详细日志信息

## 常见问题

1. API密钥未配置
   - 如果遇到API密钥相关错误，请确保已正确配置DeepSeek API密钥。

2. Git仓库问题
   - 确保在Git仓库目录下执行命令，且有足够的Git权限。

3. pre-commit钩子不生效
   - 检查钩子文件是否有执行权限
   - 确保钩子文件路径正确（.git/hooks/pre-commit）

## 更多信息

详细文档请参考[GitHub项目页面](https://github.com/icatw/ds-cr-tool)。