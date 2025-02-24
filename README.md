# GPT-CodeReview

基于 DeepSeek AI 的智能代码评审工具，通过分析 Git 差异自动生成代码评审报告，提高代码评审效率和质量。

## 特性

- 🚀 基于 DeepSeek AI 的智能代码分析
- 📊 自动生成详细的代码评审报告
- 🔄 无缝集成 Git 工作流
- ⚡ 高性能本地缓存
- 🛠️ 简单的配置和使用方式
- 🔌 可扩展的架构设计

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
  "model": "deepseek-chat",
  "cache": {
    "enabled": true,
    "ttl": 86400
  },
  "output": {
    "format": "markdown",
    "template": "default"
  }
}
```

## 贡献

欢迎提交 Pull Request 和 Issue！

1. Fork 本仓库
2. 创建您的特性分支：`git checkout -b feature/amazing-feature`
3. 提交您的改动：`git commit -m 'Add some amazing feature'`
4. 推送到分支：`git push origin feature/amazing-feature`
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情
