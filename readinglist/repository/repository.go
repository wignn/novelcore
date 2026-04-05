package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wignn/micro-3/readinglist/model"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type ReadingListRepository interface {
	Close()
	AddEntry(c context.Context, e *model.ReadingListEntry) error
	UpdateEntry(c context.Context, e *model.ReadingListEntry) error
	GetEntries(c context.Context, accountID, status string, skip, take uint64) ([]*model.ReadingListEntry, error)
	GetEntryByID(c context.Context, id string) (*model.ReadingListEntry, error)
	DeleteEntry(c context.Context, id string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		panic(err)
	}
}

func (r *PostgresRepository) AddEntry(c context.Context, e *model.ReadingListEntry) error {
	_, err := r.db.ExecContext(c,
		`INSERT INTO reading_lists (id, account_id, novel_id, status, current_chapter, rating, notes, is_favorite)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 ON CONFLICT (account_id, novel_id) DO UPDATE SET
		   status = EXCLUDED.status, current_chapter = EXCLUDED.current_chapter,
		   rating = EXCLUDED.rating, notes = EXCLUDED.notes, is_favorite = EXCLUDED.is_favorite`,
		e.ID, e.AccountID, e.NovelID, e.Status, e.CurrentChapter, e.Rating, e.Notes, e.IsFavorite)
	return err
}

func (r *PostgresRepository) UpdateEntry(c context.Context, e *model.ReadingListEntry) error {
	_, err := r.db.ExecContext(c,
		`UPDATE reading_lists SET status=$1, current_chapter=$2, rating=$3, notes=$4, is_favorite=$5 WHERE id=$6`,
		e.Status, e.CurrentChapter, e.Rating, e.Notes, e.IsFavorite, e.ID)
	return err
}

func (r *PostgresRepository) GetEntries(c context.Context, accountID, status string, skip, take uint64) ([]*model.ReadingListEntry, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT id, account_id, novel_id, status, current_chapter, rating, notes, is_favorite, created_at, updated_at
		         FROM reading_lists WHERE account_id = $1 AND status = $2 ORDER BY updated_at DESC OFFSET $3 LIMIT $4`
		args = []interface{}{accountID, status, skip, take}
	} else {
		query = `SELECT id, account_id, novel_id, status, current_chapter, rating, notes, is_favorite, created_at, updated_at
		         FROM reading_lists WHERE account_id = $1 ORDER BY updated_at DESC OFFSET $2 LIMIT $3`
		args = []interface{}{accountID, skip, take}
	}

	rows, err := r.db.QueryContext(c, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*model.ReadingListEntry
	for rows.Next() {
		e := &model.ReadingListEntry{}
		var rating sql.NullInt32
		if err := rows.Scan(&e.ID, &e.AccountID, &e.NovelID, &e.Status, &e.CurrentChapter,
			&rating, &e.Notes, &e.IsFavorite, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		if rating.Valid {
			e.Rating = rating.Int32
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *PostgresRepository) GetEntryByID(c context.Context, id string) (*model.ReadingListEntry, error) {
	e := &model.ReadingListEntry{}
	var rating sql.NullInt32
	err := r.db.QueryRowContext(c,
		`SELECT id, account_id, novel_id, status, current_chapter, rating, notes, is_favorite, created_at, updated_at
		 FROM reading_lists WHERE id = $1`, id).
		Scan(&e.ID, &e.AccountID, &e.NovelID, &e.Status, &e.CurrentChapter,
			&rating, &e.Notes, &e.IsFavorite, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if rating.Valid {
		e.Rating = rating.Int32
	}
	return e, nil
}

func (r *PostgresRepository) DeleteEntry(c context.Context, id string) error {
	res, err := r.db.ExecContext(c, "DELETE FROM reading_lists WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("entry not found with id %s", id)
	}
	return nil
}
