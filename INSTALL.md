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

## 使用方法

### 基本命令

1. 评审最新的代码改动：

```bash
cr diff
```

2. 评审指定的提交：

```bash
cr review <commit-id>
```

3. 评审指定范围的提交：

```bash
cr review <start-commit>..<end-commit>
```

### 命令参数

- `--files`: 指定要评审的文件列表，多个文件用逗号分隔
- `--commit-range`: 指定要评审的提交范围
- `--format`: 输出格式（markdown/html/pdf），默认为markdown
- `--output`: 输出文件路径，默认输出到标准输出
- `--verbose`: 显示详细日志信息

### 示例

1. 评审特定文件：

```bash
cr diff --files=main.go,utils.go
```

2. 输出HTML格式报告：

```bash
cr diff --format=html --output=review.html
```

3. 评审最近3次提交：

```bash
cr review HEAD~3..HEAD
```

## 常见问题

1. API密钥未配置

如果遇到API密钥相关错误，请确保已正确配置DeepSeek API密钥。

2. Git仓库问题

确保在Git仓库目录下执行命令，且有足够的Git权限。

## 更多信息

详细文档请参考[GitHub项目页面](https://github.com/icatw/ds-cr-tool)。