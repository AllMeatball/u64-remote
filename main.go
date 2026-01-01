package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"os"
	"u64-remote/server"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("u64-remote [prg-file]")
		os.Exit(1)
	}

	prg_file := args[0]

	creds := server.U64Creds{}
	cfg_dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Config not found using working directory: %s\n", cfg_dir)
		cfg_dir = "./"
	} else {
		cfg_dir = filepath.Join(cfg_dir, "u64-remote")
	}

	creds_path := filepath.Join(cfg_dir, "creds.json")

	file, err := os.Open(creds_path)
	if err != nil { panic(err) }
	defer file.Close()

	err = json.NewDecoder(file).Decode(&creds)
	if err != nil { panic(err) }

	server, err := server.NewU64Server(creds)
	if err != nil { panic(err) }

	prg, err := os.Open(prg_file)
	if err != nil { panic(err) }
	defer prg.Close()

	server.RunPRG(prg)
}
