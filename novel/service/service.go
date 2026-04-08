package service

import (
	"context"
	"errors"
	"strings"
	"unicode"

	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/novel/model"
	"github.com/wignn/micro-3/novel/repository"
)

var ErrInvalidTagName = errors.New("tag name is required")

type NovelService interface {
	// Novel
	CreateNovel(c context.Context, title, altTitle, description, coverURL, authorID, status, novelType, country string, year int32, genreIDs, tagIDs []int32) (*model.Novel, error)
	GetNovel(c context.Context, id string) (*model.Novel, error)
	ListNovels(c context.Context, skip, take uint64, query, status, novelType, country, sortBy, sortOrder string, genreIDs, tagIDs []int32) ([]*model.Novel, error)
	UpdateNovel(c context.Context, id, title, altTitle, description, coverURL, authorID, status, novelType, country string, year int32, genreIDs, tagIDs []int32) (*model.Novel, error)
	DeleteNovel(c context.Context, id string) error

	// Chapter
	CreateChapter(c context.Context, novelID string, chapterNumber float64, title, translatorGroupID, sourceURL string) (*model.Chapter, error)
	GetChapter(c context.Context, id string) (*model.Chapter, error)
	ListChapters(c context.Context, novelID string, skip, take uint64) ([]*model.Chapter, error)
	UpdateChapter(c context.Context, id, novelID string, chapterNumber float64, title, translatorGroupID, sourceURL string) (*model.Chapter, error)
	DeleteChapter(c context.Context, id string) error

	// Author
	CreateAuthor(c context.Context, name, bio string) (*model.Author, error)
	ListAuthors(c context.Context, skip, take uint64, id string) ([]*model.Author, error)

	// Translation Group
	CreateTranslationGroup(c context.Context, name, websiteURL, description string) (*model.TranslationGroup, error)
	ListTranslationGroups(c context.Context, skip, take uint64) ([]*model.TranslationGroup, error)

	// Genre & Tag
	GetGenres(c context.Context) ([]model.Genre, error)
	CreateTag(c context.Context, name, slug string) (*model.Tag, error)
	GetTags(c context.Context) ([]model.Tag, error)

	// Ranking
	GetRanking(c context.Context, period, sortBy string, skip, take uint64) ([]*model.Novel, error)

	// View
	IncrementViewCount(c context.Context, novelID string) (int64, error)
}

type novelService struct {
	repo repository.NovelRepository
}

func NewNovelService(r repository.NovelRepository) NovelService {
	return &novelService{repo: r}
}

func (s *novelService) CreateNovel(c context.Context, title, altTitle, description, coverURL, authorID, status, novelType, country string, year int32, genreIDs, tagIDs []int32) (*model.Novel, error) {
	if status == "" {
		status = "ongoing"
	}
	if novelType == "" {
		novelType = "web_novel"
	}

	n := &model.Novel{
		ID:               ksuid.New().String(),
		Title:            title,
		AlternativeTitle: altTitle,
		Description:      description,
		CoverImageURL:    coverURL,
		AuthorID:         authorID,
		Status:           status,
		NovelType:        novelType,
		CountryOfOrigin:  country,
		YearPublished:    year,
	}

	if err := s.repo.CreateNovel(c, n, genreIDs, tagIDs); err != nil {
		return nil, err
	}

	return s.repo.GetNovelByID(c, n.ID)
}

func (s *novelService) GetNovel(c context.Context, id string) (*model.Novel, error) {
	return s.repo.GetNovelByID(c, id)
}

func (s *novelService) ListNovels(c context.Context, skip, take uint64, query, status, novelType, country, sortBy, sortOrder string, genreIDs, tagIDs []int32) ([]*model.Novel, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}

	if query != "" {
		return s.repo.SearchNovels(c, query, skip, take)
	}

	return s.repo.ListNovels(c, skip, take, status, novelType, country, sortBy, sortOrder, genreIDs, tagIDs)
}

func (s *novelService) UpdateNovel(c context.Context, id, title, altTitle, description, coverURL, authorID, status, novelType, country string, year int32, genreIDs, tagIDs []int32) (*model.Novel, error) {
	n := &model.Novel{
		ID:               id,
		Title:            title,
		AlternativeTitle: altTitle,
		Description:      description,
		CoverImageURL:    coverURL,
		AuthorID:         authorID,
		Status:           status,
		NovelType:        novelType,
		CountryOfOrigin:  country,
		YearPublished:    year,
	}

	if err := s.repo.UpdateNovel(c, n, genreIDs, tagIDs); err != nil {
		return nil, err
	}

	return s.repo.GetNovelByID(c, id)
}

func (s *novelService) DeleteNovel(c context.Context, id string) error {
	return s.repo.DeleteNovel(c, id)
}

func (s *novelService) CreateChapter(c context.Context, novelID string, chapterNumber float64, title, translatorGroupID, sourceURL string) (*model.Chapter, error) {
	ch := &model.Chapter{
		ID:                ksuid.New().String(),
		NovelID:           novelID,
		ChapterNumber:     chapterNumber,
		Title:             title,
		TranslatorGroupID: translatorGroupID,
		SourceURL:         sourceURL,
	}

	if err := s.repo.CreateChapter(c, ch); err != nil {
		return nil, err
	}

	return s.repo.GetChapterByID(c, ch.ID)
}

func (s *novelService) GetChapter(c context.Context, id string) (*model.Chapter, error) {
	return s.repo.GetChapterByID(c, id)
}

func (s *novelService) ListChapters(c context.Context, novelID string, skip, take uint64) ([]*model.Chapter, error) {
	if take > 500 || (take == 0 && skip == 0) {
		take = 500
	}
	return s.repo.ListChapters(c, novelID, skip, take)
}

func (s *novelService) UpdateChapter(c context.Context, id, novelID string, chapterNumber float64, title, translatorGroupID, sourceURL string) (*model.Chapter, error) {
	ch := &model.Chapter{
		ID:                id,
		NovelID:           novelID,
		ChapterNumber:     chapterNumber,
		Title:             title,
		TranslatorGroupID: translatorGroupID,
		SourceURL:         sourceURL,
	}

	if err := s.repo.UpdateChapter(c, ch); err != nil {
		return nil, err
	}

	return s.repo.GetChapterByID(c, id)
}

func (s *novelService) DeleteChapter(c context.Context, id string) error {
	return s.repo.DeleteChapter(c, id)
}

func (s *novelService) CreateAuthor(c context.Context, name, bio string) (*model.Author, error) {
	a := &model.Author{
		ID:   ksuid.New().String(),
		Name: name,
		Bio:  bio,
	}

	if err := s.repo.CreateAuthor(c, a); err != nil {
		return nil, err
	}

	return s.repo.GetAuthorByID(c, a.ID)
}

func (s *novelService) ListAuthors(c context.Context, skip, take uint64, id string) ([]*model.Author, error) {
	if id != "" {
		a, err := s.repo.GetAuthorByID(c, id)
		if err != nil {
			return nil, err
		}
		return []*model.Author{a}, nil
	}

	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}
	return s.repo.ListAuthors(c, skip, take)
}

func (s *novelService) CreateTranslationGroup(c context.Context, name, websiteURL, description string) (*model.TranslationGroup, error) {
	g := &model.TranslationGroup{
		ID:          ksuid.New().String(),
		Name:        name,
		WebsiteURL:  websiteURL,
		Description: description,
	}

	if err := s.repo.CreateTranslationGroup(c, g); err != nil {
		return nil, err
	}

	return g, nil
}

func (s *novelService) ListTranslationGroups(c context.Context, skip, take uint64) ([]*model.TranslationGroup, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}
	return s.repo.ListTranslationGroups(c, skip, take)
}

func (s *novelService) GetGenres(c context.Context) ([]model.Genre, error) {
	return s.repo.GetGenres(c)
}

func (s *novelService) CreateTag(c context.Context, name, slug string) (*model.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidTagName
	}

	slug = strings.TrimSpace(slug)
	if slug == "" {
		slug = toSlug(name)
	}

	return s.repo.CreateTag(c, name, slug)
}

func (s *novelService) GetTags(c context.Context) ([]model.Tag, error) {
	return s.repo.GetTags(c)
}

func toSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			prevDash = false
		case !prevDash:
			b.WriteByte('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func (s *novelService) GetRanking(c context.Context, period, sortBy string, skip, take uint64) ([]*model.Novel, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 50
	}
	return s.repo.GetRanking(c, period, sortBy, skip, take)
}

func (s *novelService) IncrementViewCount(c context.Context, novelID string) (int64, error) {
	return s.repo.IncrementViewCount(c, novelID)
}
