package pgsql

import (
	"context"
	"fmt"
	"url_shortener/pkg/urlGenerator"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool    *pgxpool.Pool
	urlsize uint
}

func New(pool *pgxpool.Pool, urlsize uint) *DB {
	return &DB{pool: pool, urlsize: urlsize}
}

func (db *DB) URLSize() uint {
	return db.urlsize
}
func (db *DB) MakeShort(url string) (string, error) {
	short, err := db.getShort(url)
	if err != nil {
		fmt.Println("get short")
		return "", err
	}

	if short != "" {
		return short, nil
	}
	short, err = db.addShort(url)
	if err != nil {
		fmt.Println("add short")
		return "", err
	}
	return short, nil
}

func (db *DB) addShort(url string) (string, error) {
	short := urlGenerator.Generate(db.urlsize)
	_, err := db.pool.Exec(context.Background(), `INSERT INTO urls (original, short) VALUES ($1, $2)`, url, short)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				switch pgErr.ConstraintName {
				case "original": //race
					return db.getShort(url)
				case "short": //collision
					return db.addShort(url)
				default:
					return "", err
				}
			}
		}
		return "", err
	}
	return short, nil
}

func (db *DB) getShort(url string) (string, error) {
	ctx := context.Background()
	row := db.pool.QueryRow(ctx, `SELECT short FROM urls WHERE original = $1`, url)
	short := ""
	err := row.Scan(&short)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return short, nil
}

func (db *DB) GetOriginal(url string) (string, error) {
	ctx := context.Background()
	row := db.pool.QueryRow(ctx, `SELECT original FROM urls WHERE short = $1`, url)

	original := ""
	err := row.Scan(&original)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return original, nil
}
