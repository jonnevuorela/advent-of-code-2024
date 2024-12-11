package ui

import (
	"syscall"
	"unsafe"
)

/**
* @brief Gets terminal size by columns and rows
*
* @param
* @param int
*
* @return
 */
func GetTerminalSize() (int, int) {
	var ws struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		return 24, 80 // Default size
	}
	return int(ws.Col), int(ws.Row)
}
