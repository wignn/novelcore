package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wignn/micro-3/review/model"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type ReviewRepository interface {
	Close()
	PutReview(c context.Context, r *model.Review) error
	GetReviewById(c context.Context, id string) (*model.Review, error)
	GetReviewsByNovel(c context.Context, novelID string, skip, take uint64) ([]*model.Review, error)
	DeleteReview(c context.Context, id string) error
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

func (r *PostgresRepository) PutReview(c context.Context, rev *model.Review) error {
	_, err := r.db.ExecContext(c,
		`INSERT INTO reviews (id, novel_id, account_id, rating, title, content, is_spoiler, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (account_id, novel_id) DO UPDATE SET
		   rating = EXCLUDED.rating, title = EXCLUDED.title, content = EXCLUDED.content,
		   is_spoiler = EXCLUDED.is_spoiler`,
		rev.ID, rev.NovelID, rev.AccountID, rev.Rating, rev.Title, rev.Content, rev.IsSpoiler, rev.CreatedAt)
	return err
}

func (r *PostgresRepository) GetReviewById(c context.Context, id string) (*model.Review, error) {
	row := r.db.QueryRowContext(c,
		`SELECT id, novel_id, account_id, rating, title, content, is_spoiler, upvotes, downvotes, created_at
		 FROM reviews WHERE id = $1`, id)
	rev := &model.Review{}
	if err := row.Scan(&rev.ID, &rev.NovelID, &rev.AccountID, &rev.Rating, &rev.Title,
		&rev.Content, &rev.IsSpoiler, &rev.Upvotes, &rev.Downvotes, &rev.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rev, nil
}

func (r *PostgresRepository) GetReviewsByNovel(c context.Context, novelID string, skip, take uint64) ([]*model.Review, error) {
	rows, err := r.db.QueryContext(c,
		`SELECT id, novel_id, account_id, rating, title, content, is_spoiler, upvotes, downvotes, created_at
		 FROM reviews WHERE novel_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		novelID, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		rev := &model.Review{}
		if err := rows.Scan(&rev.ID, &rev.NovelID, &rev.AccountID, &rev.Rating, &rev.Title,
			&rev.Content, &rev.IsSpoiler, &rev.Upvotes, &rev.Downvotes, &rev.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}

	return reviews, nil
}

func (r *PostgresRepository) DeleteReview(c context.Context, id string) error {
	res, err := r.db.ExecContext(c, "DELETE FROM reviews WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("review not found with id %s", id)
	}
	return nil
}
