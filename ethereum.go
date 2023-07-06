package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

func generateEthereumKeys(numAddresses int, numWorkers int) map[string]string {
	addresses := make(map[string]string)
	keys := make(chan int, numWorkers)
	results := make(chan struct {
		index   int
		pubKey  string
		privKey string
	}, numAddresses)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Worker function
	worker := func() {
		defer wg.Done()
		for index := range keys {
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}

			pubKey := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
			privKey := fmt.Sprintf("%x", privateKey.D)

			results <- struct {
				index   int
				pubKey  string
				privKey string
			}{index: index, pubKey: strings.ToLower(pubKey), privKey: privKey}
		}
	}

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Generate addresses
	go func() {
		for i := 0; i < numAddresses; i++ {
			keys <- i
		}
		close(keys)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	for result := range results {
		addresses[result.pubKey] = result.privKey
	}

	return addresses
}
