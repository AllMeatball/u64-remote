package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"github.com/AllMeatball/u64-remote/server"
)

func main() {
	creds := server.U64Creds{
		EnableMessageBox: true,
	}
	cfg_dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Config dir not found using working directory: %s\n", cfg_dir)
		cfg_dir = "./"
	} else {
		cfg_dir = filepath.Join(cfg_dir, "u64-remote")
	}

	creds_path := filepath.Join(cfg_dir, "creds.json")

	file, err := os.Open(creds_path)
	if err != nil { errHand_Fatal(err, creds.EnableMessageBox) }
	defer file.Close()

	err = json.NewDecoder(file).Decode(&creds)
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
