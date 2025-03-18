package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var (
	databaseName = "sqlite-database.db"
)

func createDb() {
	log.Println("Creating " + databaseName)
	file, err := os.Create(databaseName)
	if err != nil {
		log.Fatalf(err.Error())
	}
	file.Close()
	log.Println(databaseName + " Created")

}

func deleteDb() {
	os.Remove(databaseName)
}

func initDb() {
	sqliteDatabase, _ := sql.Open("sqlite3", "./"+databaseName)
	defer sqliteDatabase.Close()

	createCommitmentTable(sqliteDatabase)
}

type Commitments struct {
	ID             int64
	Secret         string
	Commitment     string
	Unlock_address string
	Unlock_amount  float64
	Valid_from     uint64
	Valid_till     uint64
}

func createCommitmentTable(db *sql.DB) {
	createCommitmentTableSQL := `CREATE TABLE commitments(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"secret" TEXT,
		"commitment" TEXT,
		"unlock_address" TEXT,
		"unlock_amount" REAL,
		"valid_from" INTEGER,
		"valid_till" INTEGER
		`

	log.Println("Creating commitment table....")
	statement, err := db.Prepare(createCommitmentTableSQL)
	if err != nil {
		log.Fatal((err.Error()))
	}
	statement.Exec()
	log.Println("commitment table created")
}

//TODO: instead of storing payment proofs, it should store payment proof hashes
//TODO: The payment proof hashes are used for nullification, so we can see if a payment proof was used already

func insertCommitment(
	db *sql.DB,
	secret string,
	commitment string,
	unlock_address string,
	unlock_amount float64,
	valid_from uint64,
	valid_till uint64,
) {
	log.Println("Inserting commitment record ....")
	insertCommitmentSQL := `INSERT INTO commitments(
		secret, 
		commitment, 
		unlock_address, 
		unlock_amount, 
		valid_from, 
		valid_till) VALUES (?, ?, ?, ?, ?, ?)`

	statement, err := db.Prepare(insertCommitmentSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(secret, commitment, unlock_address, unlock_amount, valid_from, valid_till)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getCommitmentDetails(db *sql.DB, commitment string) (Commitments, error) {
	var commitments Commitments

	row := db.QueryRow("SELECT unlock_address,unlock_amount,valid_from,valid_till FROM commitments WHERE commitment = ?", commitment)

	if err := row.Scan(&commitments.Unlock_address, &commitments.Unlock_amount, &commitments.Valid_from, &commitments.Valid_till); err != nil {

		if err == sql.ErrNoRows {
			return commitments, fmt.Errorf("commitment: %x: no such commitment", commitment)
		}
		return commitments, fmt.Errorf("commtiment: %x: %v", commitment, err)
	}

	return commitments, nil
}

func getSecret(db *sql.DB, commitment string) (Commitments, error) {
	var commitments Commitments

	row := db.QueryRow("SELECT secret, unlock_address,unlock_amount,valid_from,valid_till FROM commitments WHERE commitment = ?", commitment)

	if err := row.Scan(&commitments.Secret, &commitments.Unlock_address, &commitments.Unlock_amount, &commitments.Valid_from, &commitments.Valid_till); err != nil {

		if err == sql.ErrNoRows {
			return commitments, fmt.Errorf("commitment: %x: no such commitment", commitment)
		}
		return commitments, fmt.Errorf("commtiment: %x: %v", commitment, err)
	}

	return commitments, nil
}
