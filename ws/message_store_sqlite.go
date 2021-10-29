package ws

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
type SqliteMessageStore struct {
	Db *sql.DB
}

func (p *SqliteMessageStore) Initialize() error {
	db, err := sql.Open("sqlite3", "wsManager.db")
	if err != nil {
		log.Printf("Failed to open database, err %s", err.Error())
		return err
	}
	p.Db = db
	// Create database if not already exists.
	stmt, err := p.Db.Prepare("CREATE TABLE IF NOT EXISTS listing_events (id INTEGER PRIMARY KEY, " +
		"user_id TEXT, message TEXT, received_at timestamp, sent_at int)")
	if err != nil {
		log.Printf("Failed to create table, store intialization failed, err %s", err.Error())
		return err
	}
	stmt.Exec()
	return nil
}

func (p *SqliteMessageStore) StoreMessage(userID , message string, receivedAt time.Time) (int64, error) {
	// We want to capture the last insert id so serializing here.
	mutex.Lock()
	id := 0
	stmt, err := p.Db.Prepare("INSERT INTO listing_events(user_id, message, received_at, sent_at) " +
		"VALUES (?, ?, ?, 0)")
	if err != nil {
		log.Printf("Failed to insert listing events into db: err %s", err.Error())
		mutex.Unlock()
		return int64(id), err
	}
	stmt.Exec(userID, message, receivedAt)
	rows, err := p.Db.Query("select last_insert_rowid()")
	var str string
	// should only be one
	for rows.Next() {
		rows.Scan(&str)
	}
	id, _ = strconv.Atoi(str)
	mutex.Unlock()
	return int64(id), nil
}

func (p *SqliteMessageStore) RetrieveUnsentMessages(userID string) []DbMessage {
	rows, _ := p.Db.Query("SELECT id, user_id, message, received_at, sent_at FROM listing_events WHERE " +
		"user_id = $1 and sent_at = 0", userID)
	dbMessage := make([]DbMessage, 0)
	for rows.Next() {
		var id int64
		var userId string
		var message string
		var receivedAt time.Time
		var sent int
		var sent_at bool
		rows.Scan(&id, &userId, &message, &receivedAt, &sent)
		if sent == 0 {
			sent_at = false
		} else {
			sent_at = true
		}
		tempMessage := DbMessage {
			ID: &id,
			UserID: &userId,
			Message: &message,
			ReceivedAt: &receivedAt,
			Sent: &sent_at,
		}
		dbMessage = append(dbMessage, tempMessage)
	}
	log.Printf("Found %d rows of unsent events for user %s",
		len(dbMessage), userID)
	return dbMessage
}

func (p *SqliteMessageStore) ConfirmSentMessage(id int64)  error {
	// mark this event as confirm sent
	stmt, err := p.Db.Prepare("UPDATE listing_events set sent_at = 1 where id = ?")
	if err != nil {
		log.Printf("Failed to mark event confirmed sent. Event id : %d, err: %s", id,
			err.Error())
		return err
	}
	stmt.Exec(id)
	return nil
}

func (p *SqliteMessageStore) Close() {
	p.Db.Close()
}

