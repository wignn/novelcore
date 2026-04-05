package model

import "time"

type Novel struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	AlternativeTitle string    `json:"alternative_title"`
	Description      string    `json:"description"`
	CoverImageURL    string    `json:"cover_image_url"`
	AuthorID         string    `json:"author_id"`
	Status           string    `json:"status"`
	NovelType        string    `json:"novel_type"`
	CountryOfOrigin  string    `json:"country_of_origin"`
	YearPublished    int32     `json:"year_published"`
	TotalChapters    int32     `json:"total_chapters"`
	RatingAvg        float64   `json:"rating_avg"`
	RatingCount      int32     `json:"rating_count"`
	ViewCount        int64     `json:"view_count"`
	BookmarkCount    int32     `json:"bookmark_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Genres           []Genre   `json:"genres,omitempty"`
	Tags             []Tag     `json:"tags,omitempty"`
	Author           *Author   `json:"author,omitempty"`
}

type Chapter struct {
	ID                 string    `json:"id"`
	NovelID            string    `json:"novel_id"`
	ChapterNumber      float64   `json:"chapter_number"`
	Title              string    `json:"title"`
	TranslatorGroupID  string    `json:"translator_group_id"`
	SourceURL          string    `json:"source_url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	TranslatorGroup    *TranslationGroup `json:"translator_group,omitempty"`
}

type Author struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type TranslationGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	WebsiteURL  string    `json:"website_url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Genre struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Tag struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type NovelDailyStat struct {
	ID             int32     `json:"id"`
	NovelID        string    `json:"novel_id"`
	StatDate       time.Time `json:"stat_date"`
	Views          int32     `json:"views"`
	BookmarksAdded int32     `json:"bookmarks_added"`
	ReviewsAdded   int32     `json:"reviews_added"`
}

type RankedNovel struct {
	Rank   int32   `json:"rank"`
	Novel  *Novel  `json:"novel"`
	Score  float64 `json:"score"`
	Change int32   `json:"change"`
}

// NovelDocument is used for Elasticsearch indexing
type NovelDocument struct {
	Title            string  `json:"title"`
	AlternativeTitle string  `json:"alternative_title"`
	Description      string  `json:"description"`
	AuthorName       string  `json:"author_name"`
	Status           string  `json:"status"`
	NovelType        string  `json:"novel_type"`
	CountryOfOrigin  string  `json:"country_of_origin"`
	Genres           []string `json:"genres"`
	Tags             []string `json:"tags"`
	RatingAvg        float64 `json:"rating_avg"`
	ViewCount        int64   `json:"view_count"`
	BookmarkCount    int32   `json:"bookmark_count"`
}
