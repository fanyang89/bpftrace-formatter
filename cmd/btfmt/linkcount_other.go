//go:build !unix

package main

import "os"

func fileLinkCount(_ os.FileInfo) (uint64, bool) {
	return 0, false
}
