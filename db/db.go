package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	dbase sql.DB
)

func InitializeDb() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	dbase = *db
	if err != nil {
		log.Fatalf("Could not open a connection the the database:\n", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS guilds (guildID TEXT not null primary key, action TEXT, logChannel TEXT)`)
	if err != nil {
		log.Fatalf("Could not create database table:\n%s", err)
	}
}
func GetServerOption(guildID, option string, passTo interface{}) error {
	row := dbase.QueryRow(fmt.Sprintf("SELECT %s FROM guilds WHERE guildID = ?;", option), guildID)
	return row.Scan(passTo)
}
func SetServerOption(guildID, option, optionValue string) error {
	_, err := dbase.Exec(fmt.Sprintf("REPLACE INTO guilds (guildID, %s) VALUES(?, ?);", option), guildID, optionValue)
	return err
}
