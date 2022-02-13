package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

var (
	dbase sql.DB
)

func InitializeDb() {
	var db *sql.DB
	var err error
	if os.Getenv("BUILD") == "PROD" {
		username := os.Getenv("MYSQL_USER")
		passwd := os.Getenv("MYSQL_PASSWORD")

		// reconnection loop
		for i := 0; i < 5; i++ {

			db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(db:3306)/db", username, passwd))
			if err != nil {
				log.Fatalf("Could not open a connection the the database:\n", err)
			}
			if err = db.Ping(); err == nil {
				break
			}
			log.Println("Could not connect to database.. Retrying in 30 seconds..")
			time.Sleep(30 * time.Second)
		}

	} else {
		db, err = sql.Open("sqlite3", "db.sqlite")
		if err != nil {
			log.Fatalf("Could not open a connection the the database:\n", err)
		}
	}
	dbase = *db

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS guilds (guildID VARCHAR(255) not null primary key, action TEXT, logChannel TEXT, timeoutDuration TEXT)`)
	if err != nil {
		log.Fatalf("Could not create database table:\n%s", err)
	}
}
func GetServerOption(guildID, option string, passTo interface{}) error {
	row := dbase.QueryRow(fmt.Sprintf("SELECT %s FROM guilds WHERE guildID = ?;", option), guildID)
	return row.Scan(passTo)
}
func SetServerOption(guildID, option, optionValue string) error {
	_, err := dbase.Exec(fmt.Sprintf("INSERT INTO guilds (guildID, %s) VALUES(?, ?) ON DUPLICATE KEY UPDATE %s=VALUES(%s);", option, option, option), guildID, optionValue)
	return err
}
func RemoveServer(guildID string) error {
	_, err := dbase.Exec("DELETE FROM guilds WHERE guildID = ?", guildID)
	return err
}
