package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// BotRecipient class
type BotRecipient struct {
	uid    string
	isChat bool
}

// Recipient method
func (r BotRecipient) Recipient() string {
	return r.uid
}

// RecipientBuilder method
func RecipientBuilder(id string, isChat bool) BotRecipient {
	return BotRecipient{uid: id, isChat: isChat}
}

// DB class
type DB struct {
	db *sql.DB
}

// IsSubscribed method
func (d *DB) IsSubscribed(id string) bool {
	var isSubscribed bool
	err := d.db.QueryRow("SELECT COUNT(*) FROM subscribers WHERE id = ?", id).Scan(&isSubscribed)
	if err != nil {
		log.Fatal(err)
	}
	return isSubscribed
}

// AddSubscriberRecipient method
func (d *DB) AddSubscriberRecipient(r BotRecipient) error {
	return d.AddSubscriber(r.uid, r.isChat)
}

// AddSubscriber method
func (d *DB) AddSubscriber(id string, isChat bool) error {
	tx, err := d.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO subscribers (
		id,
		isChat
	)
	VALUES (?,?);`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, isChat)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RemoveSubscriber method
func (d *DB) RemoveSubscriber(id string) error {
	tx, err := d.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`DELETE FROM subscribers WHERE id = ?`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetSubscribers method
func (d *DB) GetSubscribers() ([]BotRecipient, error) {
	var subscribers []BotRecipient
	rows, err := d.db.Query("select id, isChat from subscribers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var isChat bool
		err = rows.Scan(&id, &isChat)
		if err != nil {
			return nil, err
		}
		subscribers = append(subscribers, RecipientBuilder(id, isChat))
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

// Close method
func (d *DB) Close() error {
	return d.db.Close()
}

// DbBuilder method
func DbBuilder(dbPath string) DB {
	initizalizeDb := false

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		initizalizeDb = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if initizalizeDb == true {
		sqlStmt := `CREATE TABLE subscribers (
			id     TEXT  NOT NULL
						   PRIMARY KEY,
			isChat BOOLEAN
		);		
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Fatalf("%q: %s\n", err, sqlStmt)
		}
	}

	return DB{db: db}
}
