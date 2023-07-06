package main

import (
	"database/sql"
	"fmt"
	"time"
)

func saveKey(chain string, publicKey string, privateKey string, source string, pattern string) error {
	// Open the SQLite database or create it if it doesn't exist
	db, err := sql.Open("sqlite3", "keygen.db")
	if err != nil {
		fmt.Printf("Failed to open database: %v", err)
		return err
	}
	defer db.Close()

	// Create the 'keys' table if it doesn't already exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS keys (
			date TEXT,
			chain TEXT,
			public_key TEXT,
			private_key TEXT,
			source TEXT,
			pattern TEXT
		);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("Failed to create keys table: %v\n", err)
		return err
	}

	// Prepare the SQL statement for inserting the record
	insertSQL := "INSERT INTO keys (date, chain, public_key, private_key, source, pattern) VALUES (?, ?, ?, ?, ?, ?);"
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		fmt.Printf("Failed to prepare keys insert statement: %v\n", err)
	}
	defer stmt.Close()

	// Get the current date and time
	date := time.Now().Format("2006-01-02 15:04:05")

	// Execute the prepared statement to insert the record
	_, err = stmt.Exec(date, chain, publicKey, privateKey, source, pattern)
	if err != nil {
		fmt.Printf("Failed to insert keys record: %v\n", err)
		return err
	}

	return nil
}

func saveRun(chain string, count int) error {
	// Open the SQLite database or create it if it doesn't exist
	db, err := sql.Open("sqlite3", "keygen.db")
	if err != nil {
		fmt.Printf("Failed to open database: %v", err)
		return err
	}
	defer db.Close()

	// Create the 'tracking' table if it doesn't already exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS tracking (
			chain TEXT,
			count INTEGER
		);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("Failed to create tracking table: %v\n", err)
		return err
	}

	// Check if a row exists where chain = chain
	var rowCount int
	row := db.QueryRow("SELECT COUNT(*) FROM tracking WHERE chain = ?", chain)
	err = row.Scan(&rowCount)
	if err != nil {
		fmt.Printf("Failed to get tracking row count: %v\n", err)
		return err
	}
	// If row does not exist, insert. If it does exist, update count.
	if rowCount == 0 {
		// Prepare the SQL statement for inserting the record
		insertSQL := "INSERT INTO tracking (chain, count) VALUES (?, ?);"
		stmt, err := db.Prepare(insertSQL)
		if err != nil {
			fmt.Printf("Failed to prepare tracking insert statement: %v\n", err)
		}
		defer stmt.Close()

		// Execute the prepared statement to insert the record
		_, err = stmt.Exec(chain, count)
		if err != nil {
			fmt.Printf("Failed to insert tracking record: %v\n", err)
			return err
		}
	} else {
		// Prepare the SQL statement for updating the record
		updateSQL := "UPDATE tracking SET count = count + ? WHERE chain = ?;"
		stmt, err := db.Prepare(updateSQL)
		if err != nil {
			fmt.Printf("Failed to prepare tracking update statement: %v\n", err)
		}
		defer stmt.Close()

		// Execute the prepared statement to update the record
		_, err = stmt.Exec(count, chain)
		if err != nil {
			fmt.Printf("Failed to update tracking record: %v\n", err)
			return err
		}
	}

	return nil
}
