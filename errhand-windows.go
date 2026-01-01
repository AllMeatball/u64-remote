//go:build windows

package main

import (
	"log"
	"syscall"
	"unsafe"
)

const MB_OK        = 0x00000000
const MB_ICONERROR = 0x00000010

func errHand_Fatal(v error, gui_error bool) {
	if gui_error {
		user32 := syscall.NewLazyDLL("user32.dll")
		messageBox := user32.NewProc("MessageBoxW")

		title, _ := syscall.UTF16PtrFromString("u64-remote Error")
		message, _ := syscall.UTF16PtrFromString(v.Error())

		messageBox.Call(
			uintptr(0),
			uintptr(unsafe.Pointer(message)),
			uintptr(unsafe.Pointer(title)),
			uintptr(MB_OK | MB_ICONERROR),
		)
	}

	log.Fatal(v)
}
