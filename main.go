package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chain := os.Args[1]
	var pageSize int

	// Check if os.Args[2] (pageSize) was provided
	if len(os.Args) >= 3 {
		pageSize, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error with page size command argument:", err)
			return
		}
	} else {
		pageSize, err = strconv.Atoi(os.Getenv("PAGE_SIZE")) // Set default page size
		if err != nil {
			fmt.Println("Error retrieving PAGE_SIZE env value:", err)
			return
		}
	}

	process(chain, pageSize)
}

func process(chain string, pageSize int) {
	fmt.Printf("\n")

	// Create key value pairs to store match patterns and addresses
	var match_patterns = getMatchPatterns(chain, "pattern")
	var match_addresses = getMatchPatterns(chain, "address")
	fmt.Printf("Found %d match patterns and %d addresses\n\n", len(match_patterns), len(match_addresses))

	if len(match_patterns)+len(match_addresses) == 0 {
		fmt.Printf("No patterns or addresses found, exiting\n")
		return
	}

	for {
		numWorkers, err := strconv.Atoi(os.Getenv("WORKER_THREADS"))
		if err != nil {
			fmt.Println("Error retrieving WORKER_THREADS env value:", err)
			return
		}

		// Generate new addresses
		var start = time.Now()
		fmt.Printf("Generating %d new %s addresses. ", pageSize, chain)
		generated_keys := make(map[string]string)

		switch chain {
		case "eth":
			generated_keys = generateEthereumKeys(pageSize, numWorkers)
		case "btc":
			generated_keys = generateBitcoinKeys(pageSize, numWorkers)
		default:
			log.Fatal("Invalid chain argument")
		}
		elapsed := time.Since(start)
		fmt.Printf("Complete: %s.\n", elapsed)

		start = time.Now()

		// Match match_patterns and match_addresses against generated_keys
		// match_patterns and match_addresses are maps with key = [public_key] and value = [source_file]
		// generated_keys is a map with key = [public_key] and value = [private_key]

		// Match to patterns, loop through addresses generated
		if len(match_patterns) > 0 {
			for generated_public_key, generated_private_key := range generated_keys {
				for pattern, pattern_source := range match_patterns {
					generated_pattern := ""
					match_pattern := ""
					switch chain {
					case "eth":
						generated_pattern = strings.ToLower(generated_public_key[2:]) // Remove the "0x" to compare
						match_pattern = strings.ToLower(pattern)
					case "btc":
						generated_pattern = strings.ToLower(generated_public_key[3:]) // Remove the "bc1" to compare
						match_pattern = strings.ToLower(pattern)
					}

					if strings.HasPrefix(generated_pattern, match_pattern) {
						fmt.Printf("Match, %s, key: %s, %s, %s\n", generated_public_key, generated_private_key, pattern_source, match_pattern)
						saveKey(chain, generated_public_key, generated_private_key, pattern_source, match_pattern)
					}
				}
			}
		}
		elapsed_pattern_scan := time.Since(start)

		// Check generated keys against match addresses
		start = time.Now()
		if len(match_addresses) > 0 {
			for address, address_source := range match_addresses {
				// Check if key exists in generated_keys
				if generated_keys[address] != "" {
					fmt.Printf("ADDRESS MATCH, %s, key: %s\n", address, generated_keys[address])
					saveKey(chain, address, generated_keys[address], address_source, address)
					return
				}
			}
		}
		elapsed_address_scan := time.Since(start)
		fmt.Printf("Scanning complete. Patterns: %s, addresses: %s\n\n", elapsed_pattern_scan, elapsed_address_scan)

		saveRun(chain, pageSize)
	}
}
