package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/wignn/micro-3/novel/model"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type NovelRepository interface {
	Close()
	// Novel
	CreateNovel(c context.Context, n *model.Novel, genreIDs, tagIDs []int32) error
	GetNovelByID(c context.Context, id string) (*model.Novel, error)
	ListNovels(c context.Context, skip, take uint64, status, novelType, country, sortBy, sortOrder string, genreIDs, tagIDs []int32) ([]*model.Novel, error)
	UpdateNovel(c context.Context, n *model.Novel, genreIDs, tagIDs []int32) error
	DeleteNovel(c context.Context, id string) error

	// Chapter
	CreateChapter(c context.Context, ch *model.Chapter) error
	GetChapterByID(c context.Context, id string) (*model.Chapter, error)
	ListChapters(c context.Context, novelID string, skip, take uint64) ([]*model.Chapter, error)
	UpdateChapter(c context.Context, ch *model.Chapter) error
	DeleteChapter(c context.Context, id string) error

	// Author
	CreateAuthor(c context.Context, a *model.Author) error
	GetAuthorByID(c context.Context, id string) (*model.Author, error)
	ListAuthors(c context.Context, skip, take uint64) ([]*model.Author, error)

	// Translation Group
	CreateTranslationGroup(c context.Context, g *model.TranslationGroup) error
	ListTranslationGroups(c context.Context, skip, take uint64) ([]*model.TranslationGroup, error)

	// Genre & Tag
	GetGenres(c context.Context) ([]model.Genre, error)
	CreateTag(c context.Context, name, slug string) (*model.Tag, error)
	GetTags(c context.Context) ([]model.Tag, error)

	// Ranking
	GetRanking(c context.Context, period, sortBy string, skip, take uint64) ([]*model.Novel, error)

	// View
	IncrementViewCount(c context.Context, novelID string) (int64, error)

	// Search (basic SQL LIKE for when ES is not available)
	SearchNovels(c context.Context, query string, skip, take uint64) ([]*model.Novel, error)
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
		log.Println("Error closing database:", err)
	}
}

func (r *PostgresRepository) CreateNovel(c context.Context, n *model.Novel, genreIDs, tagIDs []int32) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(c,
		`INSERT INTO novels (id, title, alternative_title, description, cover_image_url, author_id, status, novel_type, country_of_origin, year_published)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		n.ID, n.Title, n.AlternativeTitle, n.Description, n.CoverImageURL,
		n.AuthorID, n.Status, n.NovelType, n.CountryOfOrigin, n.YearPublished)
	if err != nil {
		return err
	}

	for _, gid := range genreIDs {
		_, err = tx.ExecContext(c, "INSERT INTO novel_genres (novel_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", n.ID, gid)
		if err != nil {
			return err
		}
	}

	for _, tid := range tagIDs {
		_, err = tx.ExecContext(c, "INSERT INTO novel_tags (novel_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", n.ID, tid)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetNovelByID(c context.Context, id string) (*model.Novel, error) {
	n := &model.Novel{}
	var authorID sql.NullString
	var yearPublished sql.NullInt32

	err := r.db.QueryRowContext(c,
		`SELECT id, title, alternative_title, description, cover_image_url, author_id,
		        status, novel_type, country_of_origin, year_published, total_chapters,
		        rating_avg, rating_count, view_count, bookmark_count, created_at, updated_at
		 FROM novels WHERE id = $1`, id).
		Scan(&n.ID, &n.Title, &n.AlternativeTitle, &n.Description, &n.CoverImageURL,
			&authorID, &n.Status, &n.NovelType, &n.CountryOfOrigin, &yearPublished,
			&n.TotalChapters, &n.RatingAvg, &n.RatingCount, &n.ViewCount, &n.BookmarkCount,
			&n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if authorID.Valid {
		n.AuthorID = authorID.String
	}
	if yearPublished.Valid {
		n.YearPublished = yearPublished.Int32
	}

	// Load genres
	genres, err := r.getNovelGenres(c, id)
	if err == nil {
		n.Genres = genres
	}

	// Load tags
	tags, err := r.getNovelTags(c, id)
	if err == nil {
		n.Tags = tags
	}

	// Load author
	if n.AuthorID != "" {
		author, err := r.GetAuthorByID(c, n.AuthorID)
		if err == nil {
			n.Author = author
		}
	}

	return n, nil
}

func (r *PostgresRepository) ListNovels(c context.Context, skip, take uint64, status, novelType, country, sortBy, sortOrder string, genreIDs, tagIDs []int32) ([]*model.Novel, error) {
	query := `SELECT DISTINCT n.id, n.title, n.alternative_title, n.description, n.cover_image_url,
	           n.author_id, n.status, n.novel_type, n.country_of_origin, n.year_published,
	           n.total_chapters, n.rating_avg, n.rating_count, n.view_count, n.bookmark_count,
	           n.created_at, n.updated_at FROM novels n`

	var conditions []string
	var args []interface{}
	argIdx := 1

	if len(genreIDs) > 0 {
		query += " JOIN novel_genres ng ON n.id = ng.novel_id"
		placeholders := make([]string, len(genreIDs))
		for i, gid := range genreIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, gid)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf("ng.genre_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(tagIDs) > 0 {
		query += " JOIN novel_tags nt ON n.id = nt.novel_id"
		placeholders := make([]string, len(tagIDs))
		for i, tid := range tagIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, tid)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf("nt.tag_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("n.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if novelType != "" {
		conditions = append(conditions, fmt.Sprintf("n.novel_type = $%d", argIdx))
		args = append(args, novelType)
		argIdx++
	}
	if country != "" {
		conditions = append(conditions, fmt.Sprintf("n.country_of_origin = $%d", argIdx))
		args = append(args, country)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Sort
	orderCol := "n.updated_at"
	switch sortBy {
	case "title":
		orderCol = "n.title"
	case "rating":
		orderCol = "n.rating_avg"
	case "views":
		orderCol = "n.view_count"
	case "bookmarks":
		orderCol = "n.bookmark_count"
	case "created":
		orderCol = "n.created_at"
	}

	order := "DESC"
	if sortOrder == "asc" {
		order = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", orderCol, order)

	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, skip, take)

	rows, err := r.db.QueryContext(c, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNovels(c, rows)
}

func (r *PostgresRepository) UpdateNovel(c context.Context, n *model.Novel, genreIDs, tagIDs []int32) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(c,
		`UPDATE novels SET title=$1, alternative_title=$2, description=$3, cover_image_url=$4,
		 author_id=$5, status=$6, novel_type=$7, country_of_origin=$8, year_published=$9
		 WHERE id=$10`,
		n.Title, n.AlternativeTitle, n.Description, n.CoverImageURL,
		n.AuthorID, n.Status, n.NovelType, n.CountryOfOrigin, n.YearPublished, n.ID)
	if err != nil {
		return err
	}

	// Replace genres
	_, _ = tx.ExecContext(c, "DELETE FROM novel_genres WHERE novel_id = $1", n.ID)
	for _, gid := range genreIDs {
		_, err = tx.ExecContext(c, "INSERT INTO novel_genres (novel_id, genre_id) VALUES ($1, $2)", n.ID, gid)
		if err != nil {
			return err
		}
	}

	// Replace tags
	_, _ = tx.ExecContext(c, "DELETE FROM novel_tags WHERE novel_id = $1", n.ID)
	for _, tid := range tagIDs {
		_, err = tx.ExecContext(c, "INSERT INTO novel_tags (novel_id, tag_id) VALUES ($1, $2)", n.ID, tid)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) DeleteNovel(c context.Context, id string) error {
	res, err := r.db.ExecContext(c, "DELETE FROM novels WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) CreateChapter(c context.Context, ch *model.Chapter) error {
	_, err := r.db.ExecContext(c,
		`INSERT INTO chapters (id, novel_id, chapter_number, title, translator_group_id, source_url)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		ch.ID, ch.NovelID, ch.ChapterNumber, ch.Title, ch.TranslatorGroupID, ch.SourceURL)
	if err != nil {
		return err
	}

	// Update total_chapters count
	_, _ = r.db.ExecContext(c,
		"UPDATE novels SET total_chapters = (SELECT COUNT(*) FROM chapters WHERE novel_id = $1) WHERE id = $1",
		ch.NovelID)

	return nil
}

func (r *PostgresRepository) GetChapterByID(c context.Context, id string) (*model.Chapter, error) {
	ch := &model.Chapter{}
	var groupID sql.NullString

	err := r.db.QueryRowContext(c,
		`SELECT c.id, c.novel_id, c.chapter_number, c.title, c.translator_group_id, c.source_url, c.created_at, c.updated_at
		 FROM chapters c WHERE c.id = $1`, id).
		Scan(&ch.ID, &ch.NovelID, &ch.ChapterNumber, &ch.Title, &groupID, &ch.SourceURL, &ch.CreatedAt, &ch.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if groupID.Valid {
		ch.TranslatorGroupID = groupID.String
	}

	return ch, nil
}

func (r *PostgresRepository) ListChapters(c context.Context, novelID string, skip, take uint64) ([]*model.Chapter, error) {
	rows, err := r.db.QueryContext(c,
		`SELECT c.id, c.novel_id, c.chapter_number, c.title, c.translator_group_id, c.source_url, c.created_at, c.updated_at
		 FROM chapters c WHERE c.novel_id = $1 ORDER BY c.chapter_number ASC OFFSET $2 LIMIT $3`,
		novelID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []*model.Chapter
	for rows.Next() {
		ch := &model.Chapter{}
		var groupID sql.NullString
		if err := rows.Scan(&ch.ID, &ch.NovelID, &ch.ChapterNumber, &ch.Title, &groupID, &ch.SourceURL, &ch.CreatedAt, &ch.UpdatedAt); err != nil {
			return nil, err
		}
		if groupID.Valid {
			ch.TranslatorGroupID = groupID.String
		}
		chapters = append(chapters, ch)
	}

	return chapters, nil
}

func (r *PostgresRepository) UpdateChapter(c context.Context, ch *model.Chapter) error {
	_, err := r.db.ExecContext(c,
		`UPDATE chapters SET chapter_number=$1, title=$2, translator_group_id=$3, source_url=$4 WHERE id=$5`,
		ch.ChapterNumber, ch.Title, ch.TranslatorGroupID, ch.SourceURL, ch.ID)
	return err
}

func (r *PostgresRepository) DeleteChapter(c context.Context, id string) error {
	// Get novel_id before delete
	var novelID string
	_ = r.db.QueryRowContext(c, "SELECT novel_id FROM chapters WHERE id = $1", id).Scan(&novelID)

	res, err := r.db.ExecContext(c, "DELETE FROM chapters WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	if novelID != "" {
		_, _ = r.db.ExecContext(c,
			"UPDATE novels SET total_chapters = (SELECT COUNT(*) FROM chapters WHERE novel_id = $1) WHERE id = $1",
			novelID)
	}

	return nil
}

func (r *PostgresRepository) CreateAuthor(c context.Context, a *model.Author) error {
	_, err := r.db.ExecContext(c,
		"INSERT INTO authors (id, name, bio) VALUES ($1, $2, $3)",
		a.ID, a.Name, a.Bio)
	return err
}

func (r *PostgresRepository) GetAuthorByID(c context.Context, id string) (*model.Author, error) {
	a := &model.Author{}
	err := r.db.QueryRowContext(c,
		"SELECT id, name, bio, created_at FROM authors WHERE id = $1", id).
		Scan(&a.ID, &a.Name, &a.Bio, &a.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *PostgresRepository) ListAuthors(c context.Context, skip, take uint64) ([]*model.Author, error) {
	rows, err := r.db.QueryContext(c,
		"SELECT id, name, bio, created_at FROM authors ORDER BY name ASC OFFSET $1 LIMIT $2",
		skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*model.Author
	for rows.Next() {
		a := &model.Author{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Bio, &a.CreatedAt); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}
	return authors, nil
}

func (r *PostgresRepository) CreateTranslationGroup(c context.Context, g *model.TranslationGroup) error {
	_, err := r.db.ExecContext(c,
		"INSERT INTO translation_groups (id, name, website_url, description) VALUES ($1, $2, $3, $4)",
		g.ID, g.Name, g.WebsiteURL, g.Description)
	return err
}

func (r *PostgresRepository) ListTranslationGroups(c context.Context, skip, take uint64) ([]*model.TranslationGroup, error) {
	rows, err := r.db.QueryContext(c,
		"SELECT id, name, website_url, description, created_at FROM translation_groups ORDER BY name ASC OFFSET $1 LIMIT $2",
		skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.TranslationGroup
	for rows.Next() {
		g := &model.TranslationGroup{}
		if err := rows.Scan(&g.ID, &g.Name, &g.WebsiteURL, &g.Description, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *PostgresRepository) GetGenres(c context.Context) ([]model.Genre, error) {
	rows, err := r.db.QueryContext(c, "SELECT id, name, slug FROM genres ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []model.Genre
	for rows.Next() {
		g := model.Genre{}
		if err := rows.Scan(&g.ID, &g.Name, &g.Slug); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *PostgresRepository) CreateTag(c context.Context, name, slug string) (*model.Tag, error) {
	t := &model.Tag{}
	err := r.db.QueryRowContext(c,
		"INSERT INTO tags (name, slug) VALUES ($1, $2) RETURNING id, name, slug",
		name, slug,
	).Scan(&t.ID, &t.Name, &t.Slug)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *PostgresRepository) GetTags(c context.Context) ([]model.Tag, error) {
	rows, err := r.db.QueryContext(c, "SELECT id, name, slug FROM tags ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		t := model.Tag{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *PostgresRepository) GetRanking(c context.Context, period, sortBy string, skip, take uint64) ([]*model.Novel, error) {
	var dateFilter string
	now := time.Now()

	switch period {
	case "weekly":
		dateFilter = fmt.Sprintf("AND ds.stat_date >= '%s'", now.AddDate(0, 0, -7).Format("2006-01-02"))
	case "monthly":
		dateFilter = fmt.Sprintf("AND ds.stat_date >= '%s'", now.AddDate(0, -1, 0).Format("2006-01-02"))
	default: // all_time
		dateFilter = ""
	}

	orderCol := "total_score"
	switch sortBy {
	case "views":
		orderCol = "total_views"
	case "bookmarks":
		orderCol = "total_bookmarks"
	case "rating":
		orderCol = "n.rating_avg"
	}

	var query string
	if period == "all_time" && sortBy == "rating" {
		query = fmt.Sprintf(`SELECT n.id, n.title, n.alternative_title, n.description, n.cover_image_url,
			n.author_id, n.status, n.novel_type, n.country_of_origin, n.year_published,
			n.total_chapters, n.rating_avg, n.rating_count, n.view_count, n.bookmark_count,
			n.created_at, n.updated_at
			FROM novels n WHERE n.rating_count > 0
			ORDER BY n.rating_avg DESC, n.rating_count DESC
			OFFSET $1 LIMIT $2`)
		rows, err := r.db.QueryContext(c, query, skip, take)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return r.scanNovels(c, rows)
	}

	if period == "all_time" {
		query = fmt.Sprintf(`SELECT n.id, n.title, n.alternative_title, n.description, n.cover_image_url,
			n.author_id, n.status, n.novel_type, n.country_of_origin, n.year_published,
			n.total_chapters, n.rating_avg, n.rating_count, n.view_count, n.bookmark_count,
			n.created_at, n.updated_at
			FROM novels n ORDER BY n.%s DESC OFFSET $1 LIMIT $2`, sortBy+"_count")
		if sortBy == "views" {
			query = fmt.Sprintf(`SELECT n.id, n.title, n.alternative_title, n.description, n.cover_image_url,
				n.author_id, n.status, n.novel_type, n.country_of_origin, n.year_published,
				n.total_chapters, n.rating_avg, n.rating_count, n.view_count, n.bookmark_count,
				n.created_at, n.updated_at
				FROM novels n ORDER BY n.view_count DESC OFFSET $1 LIMIT $2`)
		}
		rows, err := r.db.QueryContext(c, query, skip, take)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return r.scanNovels(c, rows)
	}

	query = fmt.Sprintf(`SELECT n.id, n.title, n.alternative_title, n.description, n.cover_image_url,
		n.author_id, n.status, n.novel_type, n.country_of_origin, n.year_published,
		n.total_chapters, n.rating_avg, n.rating_count, n.view_count, n.bookmark_count,
		n.created_at, n.updated_at
		FROM novels n
		LEFT JOIN (
			SELECT novel_id, SUM(views) as total_views, SUM(bookmarks_added) as total_bookmarks,
			       SUM(views + bookmarks_added * 3) as total_score
			FROM novel_daily_stats ds WHERE 1=1 %s GROUP BY novel_id
		) ds ON n.id = ds.novel_id
		ORDER BY %s DESC NULLS LAST
		OFFSET $1 LIMIT $2`, dateFilter, orderCol)

	rows, err := r.db.QueryContext(c, query, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanNovels(c, rows)
}

func (r *PostgresRepository) IncrementViewCount(c context.Context, novelID string) (int64, error) {
	_, err := r.db.ExecContext(c, "UPDATE novels SET view_count = view_count + 1 WHERE id = $1", novelID)
	if err != nil {
		return 0, err
	}

	// Upsert daily stats
	today := time.Now().Format("2006-01-02")
	_, _ = r.db.ExecContext(c,
		`INSERT INTO novel_daily_stats (novel_id, stat_date, views) VALUES ($1, $2, 1)
		 ON CONFLICT (novel_id, stat_date) DO UPDATE SET views = novel_daily_stats.views + 1`,
		novelID, today)

	var count int64
	err = r.db.QueryRowContext(c, "SELECT view_count FROM novels WHERE id = $1", novelID).Scan(&count)
	return count, err
}

func (r *PostgresRepository) SearchNovels(c context.Context, query string, skip, take uint64) ([]*model.Novel, error) {
	searchQuery := "%" + query + "%"
	rows, err := r.db.QueryContext(c,
		`SELECT id, title, alternative_title, description, cover_image_url,
		        author_id, status, novel_type, country_of_origin, year_published,
		        total_chapters, rating_avg, rating_count, view_count, bookmark_count,
		        created_at, updated_at
		 FROM novels
		 WHERE title ILIKE $1 OR alternative_title ILIKE $1 OR description ILIKE $1
		 ORDER BY rating_avg DESC
		 OFFSET $2 LIMIT $3`,
		searchQuery, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanNovels(c, rows)
}

func (r *PostgresRepository) scanNovels(c context.Context, rows *sql.Rows) ([]*model.Novel, error) {
	var novels []*model.Novel
	for rows.Next() {
		n := &model.Novel{}
		var authorID sql.NullString
		var yearPublished sql.NullInt32
		if err := rows.Scan(&n.ID, &n.Title, &n.AlternativeTitle, &n.Description, &n.CoverImageURL,
			&authorID, &n.Status, &n.NovelType, &n.CountryOfOrigin, &yearPublished,
			&n.TotalChapters, &n.RatingAvg, &n.RatingCount, &n.ViewCount, &n.BookmarkCount,
			&n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		if authorID.Valid {
			n.AuthorID = authorID.String
		}
		if yearPublished.Valid {
			n.YearPublished = yearPublished.Int32
		}

		if n.AuthorID != "" {
			author, err := r.GetAuthorByID(c, n.AuthorID)
			if err == nil {
				n.Author = author
			}
		}

		genres, _ := r.getNovelGenres(c, n.ID)
		n.Genres = genres
		tags, _ := r.getNovelTags(c, n.ID)
		n.Tags = tags

		novels = append(novels, n)
	}
	return novels, nil
}

func (r *PostgresRepository) getNovelGenres(c context.Context, novelID string) ([]model.Genre, error) {
	rows, err := r.db.QueryContext(c,
		`SELECT g.id, g.name, g.slug FROM genres g
		 JOIN novel_genres ng ON g.id = ng.genre_id
		 WHERE ng.novel_id = $1 ORDER BY g.name`, novelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []model.Genre
	for rows.Next() {
		g := model.Genre{}
		if err := rows.Scan(&g.ID, &g.Name, &g.Slug); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *PostgresRepository) getNovelTags(c context.Context, novelID string) ([]model.Tag, error) {
	rows, err := r.db.QueryContext(c,
		`SELECT t.id, t.name, t.slug FROM tags t
		 JOIN novel_tags nt ON t.id = nt.tag_id
		 WHERE nt.novel_id = $1 ORDER BY t.name`, novelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		t := model.Tag{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}
