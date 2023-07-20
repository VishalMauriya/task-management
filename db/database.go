package db

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) {
    var err error
    DB, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal("Error in opening database:", err)
    }
    
    // Create tasks table if it doesn't exist
    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title VARCHAR(100) NOT NULL,
        description VARCHAR(250),
        due_date VARCHAR(100) NOT NULL,
        status VARCHAR(50) NOT NULL
    )`)
    if err != nil {
        log.Fatal("Error creating tasks table:", err)
    }
}
