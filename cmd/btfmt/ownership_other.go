//go:build !unix

package main

import "os"

func fileOwnerMismatch(_ os.FileInfo) (bool, bool) {
	return false, false
}
