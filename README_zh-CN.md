# btfmt

[English](README.md) | [中文](README_zh-CN.md)

bpftrace 脚本格式化工具，支持 VS Code 集成。

## 功能特性

- 格式化 bpftrace 脚本，统一缩进、间距和结构
- VS Code 扩展内置二进制，安装即用
- 支持语言服务器协议 (LSP)，提供悬停文档、补全、导航和重命名
- 针对 builtin、probe target、`args` 字段、可见变量、Map 和 Macro 提供上下文补全
- 支持 imported Map 和 Macro family 的跨文件定义、引用与重命名
- 通过 JSON 配置文件自定义格式化规则
- 保留注释和 shebang

## 安装

### VS Code 扩展（推荐）

从 VS Code Marketplace 安装 [btfmt - bpftrace Language Support](https://marketplace.visualstudio.com/items?itemName=fanyang89.btfmt)，或在扩展视图中搜索 `btfmt`。

也可以从 [GitHub Releases](https://github.com/fanyang89/bpftrace-formatter/releases) 下载对应平台的 VSIX：

1. 下载对应平台的 `.vsix` 文件（如 Linux x64）
2. 在 VS Code 中按 `Ctrl+Shift+P`，运行 "Extensions: Install from VSIX..."
3. 选择下载的文件

扩展已内置 btfmt 二进制，无需额外安装 formatter。
probe target 和 `args` 字段补全使用 workspace 脚本、可移植 catalog 和当前用户可读的内核 metadata，无需 root 权限。

### CLI 二进制

从 [Releases](https://github.com/fanyang89/bpftrace-formatter/releases) 下载预编译的二进制文件：

| 平台          | 文件                        |
| ------------- | --------------------------- |
| Linux x64     | `btfmt-linux-amd64.tar.gz`  |
| macOS x64     | `btfmt-darwin-amd64.tar.gz` |
| macOS ARM64   | `btfmt-darwin-arm64.tar.gz` |
| Windows x64   | `btfmt-windows-amd64.zip`   |

Linux ARM64 和 Windows ARM64 发布产物暂缓提供。

解压并添加到 PATH：

```bash
tar -xzf btfmt-linux-amd64.tar.gz
sudo mv btfmt /usr/local/bin/
```

### 从源码构建

```bash
cargo install --git https://github.com/fanyang89/bpftrace-formatter.git
```

或克隆仓库构建：

```bash
git clone https://github.com/fanyang89/bpftrace-formatter.git
cd bpftrace-formatter
cargo build --release
./target/release/btfmt --help
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

tracepoint:syscalls:sys_enter_openat
/pid == 1234/
{
    @opens[pid] = count();
}
```

### 命令行选项

```
btfmt [options] <file.bt|-> [file2.bt ...]

选项：
  -w                     将结果写回源文件
  -i                     就地修改文件（同 -w）
  -c, -config <file>     指定配置文件路径
  -v, -verbose           启用详细输出
  --check                输入未格式化时返回非零状态
  -generate-config       生成默认配置文件
  -config-output <file>  生成配置文件的输出路径
  --force                覆盖已存在的生成配置
  -version               显示版本信息
  -help                  显示帮助信息
```

使用显式 `-` 从标准输入读取：

```bash
cat script.bt | btfmt -
btfmt --check script.bt
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
btfmt -generate-config --force  # 覆盖已有文件
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
| `spacing`     | `around_parentheses`         | false       | 括号内加空格                     |
| `spacing`     | `around_brackets`            | false       | 方括号内加空格                   |
| `spacing`     | `before_block_start`         | true        | `{` 前加空格                     |
| `line_breaks` | `empty_lines_between_probes` | 1           | 探针块之间的空行数               |
| `line_breaks` | `empty_lines_after_shebang`  | 1           | shebang 后的空行数               |
| `comments`    | `preserve_inline`            | true        | 尽量保留行内注释                 |
| `blocks`      | `brace_style`                | "next_line" | `"same_line"`、`"next_line"` 或 `"gnu"` |
| `blocks`      | `indent_statements`          | true        | 缩进块内语句                     |

## VS Code 扩展

VS Code 扩展提供：

- `.bt` 文件语法高亮
- 保存时自动格式化（在 VS Code 设置中启用）
- 格式化文档命令（`Shift+Alt+F`）
- bpftrace builtin 悬停文档
- builtin、provider、probe target、`args` 字段、关键字、可见变量、Map、Macro 参数和 imported Macro 的上下文补全
- 探针与 Macro 的文档符号
- 词法变量、Map 和 Macro family 的定义、引用、高亮与重命名
- 跨 imported `.bt` 文件的 workspace 导航与重命名

也可以直接启动语言服务器，供其他 LSP 客户端使用：

```bash
btfmt lsp
```

### 扩展设置

| 设置               | 默认值   | 说明                                     |
| ------------------ | -------- | ---------------------------------------- |
| `btfmt.serverPath` | `btfmt`  | btfmt 二进制路径（默认使用内置二进制）   |
| `btfmt.configPath` | `""`     | `.btfmt.json` 配置文件路径               |

## 许可证

[Unlicense](LICENSE)（公共领域）
