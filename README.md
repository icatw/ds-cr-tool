# GPT-CodeReview

一个简单的代码评审工具，通过获取Git差异内容并调用DeepSeek API进行智能评审，自动生成代码评审报告。

## 特性

- 🚀 基于 DeepSeek AI 的智能评审
- 🔄 自动获取Git差异内容
- 📊 生成评审报告
- 🛠️ 简单易用

## 安装

```bash
# 使用 go install 安装
go install github.com/icatw/cr@latest

# 或者从源码安装
git clone https://github.com/icatw/cr.git
cd cr
go install
```

## 快速开始

1. 配置 DeepSeek API 密钥：

```bash
# 设置环境变量
export DEEPSEEK_API_KEY=your_api_key

# 或创建配置文件
echo '{"api_key": "your_api_key"}' > ~/.cr/config.json
```

2. 在 Git 仓库中使用：

```bash
# 评审最新的代码改动
cr diff

# 评审指定的提交
cr review <commit-id>

# 评审指定范围的提交
cr review <start-commit>..<end-commit>
```

## 配置

默认配置文件位置：`~/.cr/config.json`

```json
{
  "api_key": "your_deepseek_api_key",
  "model": "deepseek-chat"
}
```

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情
