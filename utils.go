package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func formatStringToETHPublicKey(s string) string {
	// abcdefois
	t := strings.TrimSpace(s)
	t = strings.ToLower(t)
	t = strings.Replace(t, "o", "0", -1)
	t = strings.Replace(t, "i", "1", -1)
	t = strings.Replace(t, "s", "5", -1)
	if !strings.HasPrefix(t, "0x") {
		t = "0x" + t
	}
	return t
}

func getMatchPatterns(chain string, match_type string) map[string]string {
	patterns := make(map[string]string)

	matchDir := fmt.Sprintf("./match_%s_%s", chain, match_type)
	err := filepath.Walk(matchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path: %s\n", path)
			return nil // Skip current path error and continue
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".txt" && filepath.Base(path) != "readme.txt" {
			fileContents, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", path)
				return nil // Skip current file read error and continue
			}

			lines := strings.Split(string(fileContents), "\n")

			for _, line := range lines {
				line = strings.ReplaceAll(line, " ", "")

				switch chain {
				case "eth":
					if len(line) > 4 {
						line = formatStringToETHPublicKey(line)
						patterns[line] = path
					}
				case "btc":
					if len(line) > 4 {
						patterns[line] = path
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking path: %s\n", matchDir)
		return nil
	}

	return patterns
}
