package formatter

import (
	"testing"

	"github.com/fanyang89/bpftrace-formatter/config"
)

var benchInput = `#!/usr/bin/env bpftrace

// Trace file opens and closes
tracepoint:syscalls:sys_enter_openat {
    printf("openat: %s\n", str(args.filename));
}

tracepoint:syscalls:sys_enter_openat2 {
    printf("openat2: %s\n", str(args->filename));
}

// Only trace specific processes
tracepoint:syscalls:sys_enter_openat /pid == 1234/ {
    @opens[pid] = count();
}

// Count system calls by process name
tracepoint:syscalls:sys_enter_* {
    @syscalls[comm] = count();
}

// Profile CPU usage
profile:hz:99 {
    @cpu_time[cpu] = count();
}

END {
    clear(@opens);
    clear(@syscalls);
    clear(@cpu_time);
}
`

func BenchmarkFormat(b *testing.B) {
	cfg := config.DefaultConfig()
	b.ResetTimer()
	for b.Loop() {
		f := NewASTFormatter(cfg)
		_, err := f.Format(benchInput)
		if err != nil {
			b.Fatal(err)
		}
	}
}
