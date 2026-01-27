# BPFTrace Formatter

一个基于 ANTLR4 的 bpftrace 脚本格式化工具。

## 功能特性

- 基于语法解析的 bpftrace 格式化
- 可配置的缩进、间距、换行、注释、探针与代码块
- 保留 shebang 和注释
- 支持多文件与就地格式化

## 构建

```bash
go mod tidy
task build
# 或
go build ./cmd/btfmt
```

## 使用方法

### 基本用法

```bash
./btfmt <file.bt>
```

### 就地修改 / 写回

```bash
./btfmt -i <file.bt>
./btfmt -w <file1.bt> <file2.bt>
```

### 配置

```bash
./btfmt -generate-config
./btfmt -config <path/to/.btfmt.json> <file.bt>
```

选项：

- `-c`, `-config <file>`: 指定配置文件路径
- `-i`: 就地修改文件
- `-w`: 将结果写回源文件（默认输出到 stdout）
- `-v`, `-verbose`: 启用详细日志
- `-generate-config`: 生成默认配置文件
- `-config-output <file>`: 生成配置的输出路径（默认：.btfmt.json）
- `-help`: 显示帮助信息

### 示例

输入文件 (`testdata/test_input.bt`):

```bpftrace
#!/usr/bin/env bpftrace
tracepoint:syscalls:sys_enter_openat{printf("openat: %s\n",str(args.filename));}
tracepoint:syscalls:sys_enter_openat2{printf("openat2: %s\n",str(args->filename));}
tracepoint:syscalls:sys_enter_openat/pid==1234/{@opens[pid]=count();}
```

输出:

```bpftrace
#!/usr/bin/env bpftrace
tracepoint:syscalls:sys_enter_openat {
    printf("openat: %s\n",str(args.filename));
}

tracepoint:syscalls:sys_enter_openat2 {
    printf("openat2: %s\n",str(args->filename));
}

tracepoint:syscalls:sys_enter_openat/pid==1234/ {
    @opens[pid] = count();
}
```

## 配置说明

配置加载顺序如下：

1. 通过 `-config` 指定的文件
2. 当前目录或父目录中的 `.btfmt.json`
3. 家目录中的 `~/.btfmt.json`
4. 内置默认值

如果通过 `-config` 指定的文件不存在，CLI 会给出警告，并直接使用内置默认值，不再继续搜索 `.btfmt.json`。

配置使用 JSON，顶层包含以下字段：`indent`、`spacing`、
`line_breaks`、`comments`、`probes`、`blocks`。

示例：

```json
{
  "indent": { "size": 4, "use_spaces": true },
  "spacing": { "around_operators": true, "around_commas": true },
  "line_breaks": { "empty_lines_between_probes": 1, "max_line_length": 80 },
  "comments": { "preserve_inline": true },
  "probes": { "sort_probes": false },
  "blocks": { "brace_style": "next_line" }
}
```

## 语法编译

```bash
task compile-grammar
```

## 测试

```bash
task test
# 或
go test ./...
```

## 项目结构

- `cmd/btfmt/`: CLI 入口
- `formatter/`: AST 格式化逻辑与访问器
- `config/`: 配置结构与加载器
- `parser/`: 生成的 ANTLR 解析器/词法器
- `bpftrace.g4`: 语法定义
- `testdata/`: 测试脚本与样例

## 支持的语法特性

- 探针定义与谓词
- 语句（if/while/for、return、print/printf、clear/delete）
- 表达式（逻辑、算术、关系、单目）
- 内置函数与 maps
- 注释与 shebang 处理

## 开发

这个项目使用 Go 开发，欢迎贡献代码和报告问题。
