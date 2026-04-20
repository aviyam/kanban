package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() {
	exe, _ := os.Executable()
	dbPath := filepath.Join(filepath.Dir(exe), "kanban.db")

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		status INTEGER,
		title TEXT,
		description TEXT,
		position INTEGER DEFAULT 0,
		completed_at DATETIME
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// Migration: Add columns if they don't exist
	_, _ = db.Exec("ALTER TABLE tasks ADD COLUMN position INTEGER DEFAULT 0")
	_, _ = db.Exec("ALTER TABLE tasks ADD COLUMN completed_at DATETIME")
}

func SaveTask(t Task) (int, error) {
	var completedAt interface{}
	if t.completedAt != nil {
		completedAt = t.completedAt.Format(time.RFC3339)
	}

	if t.id == 0 {
		res, err := db.Exec("INSERT INTO tasks (status, title, description, position, completed_at) VALUES (?, ?, ?, ?, ?)", t.status, t.title, t.description, t.position, completedAt)
		if err != nil {
			return 0, err
		}
		id, err := res.LastInsertId()
		return int(id), err
	} else {
		_, err := db.Exec("UPDATE tasks SET status = ?, title = ?, description = ?, position = ?, completed_at = ? WHERE id = ?", t.status, t.title, t.description, t.position, completedAt, t.id)
		return t.id, err
	}
}

func DeleteTask(id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func CleanupDoneTasks(days int) error {
	if days <= 0 {
		return nil
	}
	cutoff := time.Now().AddDate(0, 0, -days).Format(time.RFC3339)
	_, err := db.Exec("DELETE FROM tasks WHERE status = ? AND completed_at < ?", done, cutoff)
	return err
}

func LoadTasks() ([]Task, error) {
	rows, err := db.Query("SELECT id, status, title, description, position, completed_at FROM tasks ORDER BY position ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var s int
		var completedAtStr sql.NullString
		err = rows.Scan(&t.id, &s, &t.title, &t.description, &t.position, &completedAtStr)
		if err != nil {
			return nil, err
		}
		t.status = status(s)
		if completedAtStr.Valid {
			parsedTime, err := time.Parse(time.RFC3339, completedAtStr.String)
			if err == nil {
				t.completedAt = &parsedTime
			}
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func CloseDB() {
	db.Close()
}
