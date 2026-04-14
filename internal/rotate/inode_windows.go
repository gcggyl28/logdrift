//go:build windows

package rotate

import "os"

// inode is not available on Windows; always return 0 so rotation is
// detected only via file size regression.
func inode(_ os.FileInfo) uint64 {
	return 0
}
