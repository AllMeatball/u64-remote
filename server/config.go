package server

import (
	"os"
	"log"
	"errors"
	"strings"
	"encoding/json"
	"path/filepath"
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

func LoadCreds(creds *U64Creds) error {
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

		return errors.New(err_msg)
	}

	return nil
}
