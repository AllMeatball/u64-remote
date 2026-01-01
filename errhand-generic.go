//go:build !windows
package main

import (
	"log"
)

func errHand_Fatal(v error, gui_error bool) {
	_ = gui_error
	log.Fatal(v)
}
