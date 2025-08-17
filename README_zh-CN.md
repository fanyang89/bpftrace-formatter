# BPFTrace Formatter

一个用于格式化 bpftrace 脚本的 Go 工具。

## 功能特性

- 自动格式化 bpftrace 探针定义
- 支持谓词（predicates）格式化
- 统一的缩进和间距
- 保留注释和 shebang
- 可配置的格式化选项

## 安装

```bash
go mod tidy
go build -o bpftrace-formatter
```

## 使用方法

### 基本用法

```bash
./bpftrace-formatter <file.bt>
```

### 示例

输入文件 (`test_input.bt`):

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
    @opens[pid]=count();
}
```

## 配置选项

可以通过代码配置以下选项：

- `SetIndentSize(int)`: 设置缩进大小（默认：4）
- `SetUseSpaces(bool)`: 使用空格而非制表符（默认：true）
- `SetProbeSpacing(int)`: 设置探针之间的空行数（默认：2）

## 测试

```bash
go test -v
```

## 项目结构

- `main.go`: 主要的格式化逻辑
- `main_test.go`: 单元测试
- `example.bt`: 格式化示例
- `test_input.bt`: 测试输入文件

## 支持的语法特性

- 探针定义（tracepoint, kprobe, uprobe 等）
- 谓词过滤
- 块语句格式化
- 注释保留
- END 块处理
- 多行语句支持

## 限制

- 目前主要处理单行探针定义
- 复杂的多行块可能需要进一步优化
- 某些特殊的 bpftrace 语法可能不被完全支持

## 开发

这个项目使用 Go 开发，欢迎贡献代码和报告问题。
