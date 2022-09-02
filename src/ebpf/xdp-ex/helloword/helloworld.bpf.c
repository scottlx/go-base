#include "common.h"

SEC("tracepoint/syscalls/sys_enter_execve")

int bpf_prog(void *ctx) {
  char msg[] = "invoke bpf_prog: Hello, World!";
  bpf_trace_printk(msg, sizeof(msg));
  return 0;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
