//go:build unix

package main

import (
	"os"
	"syscall"
)

func fileOwnerMismatch(info os.FileInfo) (bool, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return false, false
	}
	uid := uint32(os.Geteuid())
	gid := uint32(os.Getegid())
	return stat.Uid != uid || stat.Gid != gid, true
}
