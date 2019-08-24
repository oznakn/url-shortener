package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database interface
type Database interface {
	Get(short string) (string, error)
	Save(url string) (string, error)
}

type sqlite struct {
	Path string
}

func randomString() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", b[0:4])
}

func (s sqlite) Save(url string) (string, error) {
	db, err := sql.Open("sqlite3", s.Path)
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	stmt, err := tx.Prepare("insert into urls(url, short) values(?, ?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	short := randomString()
	fmt.Println(short)

	_, err = stmt.Exec(url, short)
	if err != nil {
		return "", err
	}

	tx.Commit()

	return short, nil
}

func (s sqlite) Get(short string) (string, error) {
	db, err := sql.Open("sqlite3", s.Path)
	stmt, err := db.Prepare("select url from urls where short = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var url string
	err = stmt.QueryRow(short).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s sqlite) Init() {
	c, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	sqlStmt := `create table if not exists urls (id integer not null primary key, url text, short text);`
	_, err = c.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}
