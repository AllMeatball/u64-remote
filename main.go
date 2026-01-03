package main

import (
	"os"
	"errors"
	"github.com/AllMeatball/u64-remote/server"
)

func main() {
	creds := server.U64Creds{
		EnableMessageBox: true,
	}

	err := server.LoadCreds(&creds)
	if err != nil { errHand_Fatal(err, creds.EnableMessageBox) }

	args := os.Args[1:]
	if len(args) < 1 {
		errHand_Fatal(errors.New("u64-remote [prg-file]"), creds.EnableMessageBox)
	}

	prg_file := args[0]

	server, err := server.NewU64Server(creds)
	if err != nil { errHand_Fatal(err, creds.EnableMessageBox) }

	prg, err := os.Open(prg_file)
	if err != nil { errHand_Fatal(err, creds.EnableMessageBox) }
	defer prg.Close()

	err = server.RunPRG(prg)
	if err != nil { errHand_Fatal(err, creds.EnableMessageBox) }
}
