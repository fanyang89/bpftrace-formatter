# BPFTrace Formatter

A Go tool for formatting bpftrace scripts.

## Features

- Automatic formatting of bpftrace probe definitions
- Support for predicates formatting
- Unified indentation and spacing
- Preserve comments and shebang
- Configurable formatting options

## Installation

```bash
go mod tidy
go build -o bpftrace-formatter
```

## Usage

### Basic Usage

```bash
./bpftrace-formatter <file.bt>
```

### Example

Input file (`test_input.bt`):

```bpftrace
#!/usr/bin/env bpftrace
tracepoint:syscalls:sys_enter_openat{printf("openat: %s\n",str(args.filename));}
tracepoint:syscalls:sys_enter_openat2{printf("openat2: %s\n",str(args->filename));}
tracepoint:syscalls:sys_enter_openat/pid==1234/{@opens[pid]=count();}
```

Output:

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

## Configuration Options

The following options can be configured through code:

- `SetIndentSize(int)`: Set indentation size (default: 4)
- `SetUseSpaces(bool)`: Use spaces instead of tabs (default: true)
- `SetProbeSpacing(int)`: Set number of empty lines between probes (default: 2)

## Testing

```bash
go test -v
```

## Project Structure

- `main.go`: Main formatting logic
- `main_test.go`: Unit tests
- `example.bt`: Formatting example
- `test_input.bt`: Test input file

## Supported Syntax Features

- Probe definitions (tracepoint, kprobe, uprobe, etc.)
- Predicate filtering
- Block statement formatting
- Comment preservation
- END block handling
- Multi-line statement support

## Limitations

- Currently mainly handles single-line probe definitions
- Complex multi-line blocks may need further optimization
- Some special bpftrace syntax may not be fully supported

## Development

This project is developed in Go. Contributions and issue reports are welcome.
