package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./chatroom.db")
	if err != nil {
		log.Fatal(err)
	}

	createUserTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	_, err = DB.Exec(createUserTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	createChannelTable := `CREATE TABLE IF NOT EXISTS channels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`
	_, err = DB.Exec(createChannelTable)
	if err != nil {
		log.Fatalf("Error creating channels table: %v", err)
	}

	createUserChannelTable := `CREATE TABLE IF NOT EXISTS user_channels (
		user_id INTEGER,
		channel_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(channel_id) REFERENCES channels(id),
		PRIMARY KEY(user_id, channel_id)
	);`
	_, err = DB.Exec(createUserChannelTable)
	if err != nil {
		log.Fatalf("Error creating user_channels table: %v", err)
	}
	
	createMessageTable := `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		channel_id INTEGER,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(channel_id) REFERENCES channels(id)
	);`
	_, err = DB.Exec(createMessageTable)
	if err != nil {
		log.Fatalf("Error creating messages table: %v", err)
	}

}
