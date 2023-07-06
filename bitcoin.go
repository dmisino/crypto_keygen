package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// Generate legacy bitcoin keys (starting with a "1")
func generateBitcoinKeysLegacy(numAddresses int, numWorkers int) map[string]string {
	addresses := make(map[string]string)
	keys := make(chan int, numWorkers)
	results := make(chan struct {
		index      int
		address    string
		privateKey *btcec.PrivateKey
	}, numAddresses)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Worker function
	worker := func() {
		defer wg.Done()
		for index := range keys {
			privateKey, err := btcec.NewPrivateKey(btcec.S256())
			if err != nil {
				log.Fatal(err)
			}

			publicKey := privateKey.PubKey()
			pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
			address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
			if err != nil {
				log.Fatal(err)
			}

			results <- struct {
				index      int
				address    string
				privateKey *btcec.PrivateKey
			}{index: index, address: address.EncodeAddress(), privateKey: privateKey}
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
		//index := result.index + 1
		addresses[result.address] = fmt.Sprintf("%x", result.privateKey.Serialize())
	}

	return addresses
}

// Generate segwit bitcoin keys (starting with "bc1", case insensitive)
func generateBitcoinKeys(numAddresses int, numWorkers int) map[string]string {
	addresses := make(map[string]string)
	keys := make(chan int, numWorkers)
	results := make(chan struct {
		index      int
		address    string
		privateKey *btcec.PrivateKey
	}, numAddresses)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Worker function
	worker := func() {
		defer wg.Done()
		for index := range keys {
			privateKey, err := btcec.NewPrivateKey(btcec.S256())
			if err != nil {
				log.Fatal(err)
			}

			publicKey := privateKey.PubKey()

			pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

			address, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
			if err != nil {
				log.Fatal(err)
			}

			results <- struct {
				index      int
				address    string
				privateKey *btcec.PrivateKey
			}{index: index, address: address.EncodeAddress(), privateKey: privateKey}
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
		addresses[result.address] = fmt.Sprintf("%x", result.privateKey.Serialize())
	}

	return addresses
}
