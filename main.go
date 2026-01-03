package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AllMeatball/u64-remote/server"
)


func getConfigPaths() []string {
	cfg_paths := []string{}

	U64_CRED_PATH, has_cred_envar := os.LookupEnv("U64_CFG_PATH")

	if has_cred_envar {
		u64_cred_paths := strings.SplitSeq(U64_CRED_PATH, ":")

		for path := range u64_cred_paths {
			cfg_paths = append(cfg_paths, path)
		}
	}

	cfg_paths = append(cfg_paths, "./")
	cfg_dir, err := os.UserConfigDir()

	if err != nil {
		log.Printf("Can't get config path from (os.UserConfigDir): %v\n", err)
	} else {
		cfg_dir = filepath.Join(cfg_dir, "u64-remote")
		cfg_paths = append(cfg_paths, cfg_dir)
	}

	return cfg_paths
}

func main() {
	creds := server.U64Creds{
		EnableMessageBox: true,
	}

	cfg_paths := getConfigPaths()

	config_loaded := false
	for _, cfg_dir := range cfg_paths {
		log.Printf("Trying config directory: %s\n", cfg_dir)
		creds_path := filepath.Join(cfg_dir, "creds.json")

		file, err := os.Open(creds_path)
		if err != nil { continue }
		defer file.Close()

		err = json.NewDecoder(file).Decode(&creds)
		if err != nil { continue }

		config_loaded = true
		log.Printf("Found final path: %s\n", cfg_dir)
		break
	}

	if !config_loaded {
		err_msg := "All config paths failed. valid directories are:"

		for _, cfg_dir := range cfg_paths {
			err_msg += "\n" + cfg_dir
		}

		errHand_Fatal(errors.New(err_msg), creds.EnableMessageBox)
	}

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
