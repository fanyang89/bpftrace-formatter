# btfmt

[English](README.md) | [中文](README_zh-CN.md)

bpftrace 脚本格式化工具，支持 VS Code 集成。

## 功能特性

- 格式化 bpftrace 脚本，统一缩进、间距和结构
- VS Code 扩展内置二进制，安装即用
- 支持语言服务器协议 (LSP)，可集成到各种编辑器
- 通过 JSON 配置文件自定义格式化规则
- 保留注释和 shebang

## 安装

### VS Code 扩展（推荐）

从 [Releases](https://github.com/fanyang89/bpftrace-formatter/releases) 页面下载 `btfmt-lsp` 扩展：

1. 下载对应平台的 `.vsix` 文件（如 `btfmt-lsp-0.0.2@linux-x64.vsix`）
2. 在 VS Code 中按 `Ctrl+Shift+P`，运行 "Extensions: Install from VSIX..."
3. 选择下载的文件

扩展已内置 btfmt 二进制，无需额外安装。

### CLI 二进制

从 [Releases](https://github.com/fanyang89/bpftrace-formatter/releases) 下载预编译的二进制文件：

| 平台          | 文件                        |
| ------------- | --------------------------- |
| Linux x64     | `btfmt-linux-amd64.tar.gz`  |
| Linux ARM64   | `btfmt-linux-arm64.tar.gz`  |
| macOS x64     | `btfmt-darwin-amd64.tar.gz` |
| macOS ARM64   | `btfmt-darwin-arm64.tar.gz` |
| Windows x64   | `btfmt-windows-amd64.zip`   |
| Windows ARM64 | `btfmt-windows-arm64.zip`   |

解压并添加到 PATH：

```bash
tar -xzf btfmt-linux-amd64.tar.gz
sudo mv btfmt /usr/local/bin/
```

### 从源码构建

```bash
go install github.com/fanyang89/bpftrace-formatter/cmd/btfmt@latest
```

或克隆仓库构建：

```bash
git clone https://github.com/fanyang89/bpftrace-formatter.git
cd bpftrace-formatter
go build ./cmd/btfmt
```

## 使用方法

### 格式化文件

```bash
btfmt script.bt          # 输出到 stdout
btfmt -w script.bt       # 写回文件
btfmt -w *.bt            # 格式化多个文件
```

### 示例

格式化前：

```bpftrace
#!/usr/bin/env bpftrace
tracepoint:syscalls:sys_enter_openat{printf("openat: %s\n",str(args.filename));}
tracepoint:syscalls:sys_enter_openat/pid==1234/{@opens[pid]=count();}
```

格式化后：

```bpftrace
#!/usr/bin/env bpftrace

tracepoint:syscalls:sys_enter_openat
{
    printf("openat: %s\n", str(args.filename));
}

tracepoint:syscalls:sys_enter_openat /pid == 1234/
{
    @opens[pid] = count();
}
```

### 命令行选项

```
btfmt [options] <file.bt> [file2.bt ...]

选项：
  -w                     将结果写回源文件
  -i                     就地修改文件（同 -w）
  -c, -config <file>     指定配置文件路径
  -v, -verbose           启用详细输出
  -generate-config       生成默认配置文件
  -version               显示版本信息
  -help                  显示帮助信息
```

## 配置

btfmt 按以下顺序查找配置：

1. 通过 `-config` 指定的文件
2. 当前目录或父目录中的 `.btfmt.json`
3. 家目录中的 `~/.btfmt.json`
4. 内置默认值

生成默认配置文件：

```bash
btfmt -generate-config
```

示例 `.btfmt.json`：

```json
{
  "indent": {
    "size": 4,
    "use_spaces": true
  },
  "spacing": {
    "around_operators": true,
    "around_commas": true
  },
  "line_breaks": {
    "empty_lines_between_probes": 1
  },
  "blocks": {
    "brace_style": "next_line"
  }
}
```

### 配置选项

| 分类          | 选项                         | 默认值      | 说明                             |
| ------------- | ---------------------------- | ----------- | -------------------------------- |
| `indent`      | `size`                       | 4           | 每级缩进的空格/制表符数          |
| `indent`      | `use_spaces`                 | true        | 使用空格而非制表符               |
| `spacing`     | `around_operators`           | true        | 在 `=`、`+`、`-` 等周围加空格    |
| `spacing`     | `around_commas`              | true        | 逗号后加空格                     |
| `spacing`     | `before_block_start`         | true        | `{` 前加空格                     |
| `spacing`     | `after_keywords`             | true        | `if`、`while` 等关键字后加空格   |
| `line_breaks` | `empty_lines_between_probes` | 1           | 探针块之间的空行数               |
| `line_breaks` | `empty_lines_after_shebang`  | 1           | shebang 后的空行数               |
| `blocks`      | `brace_style`                | "next_line" | `"same_line"`、`"next_line"` 或 `"gnu"` |

## VS Code 扩展

VS Code 扩展提供：

- `.bt` 文件语法高亮
- 保存时自动格式化（在 VS Code 设置中启用）
- 格式化文档命令（`Shift+Alt+F`）

### 扩展设置

| 设置               | 默认值   | 说明                                     |
| ------------------ | -------- | ---------------------------------------- |
| `btfmt.serverPath` | `btfmt`  | btfmt 二进制路径（默认使用内置二进制）   |
| `btfmt.configPath` | `""`     | `.btfmt.json` 配置文件路径               |

## 许可证

[Unlicense](LICENSE)（公共领域）
