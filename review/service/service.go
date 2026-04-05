package service

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/review/model"
	"github.com/wignn/micro-3/review/repository"
)

type ReviewService interface {
	PutReview(c context.Context, novelID, accountID string, rating int, title, content string, isSpoiler bool) (*model.Review, error)
	GetReviewById(c context.Context, id string) (*model.Review, error)
	GetReviewsByNovel(c context.Context, novelID string, skip, take uint64) ([]*model.Review, error)
	DeleteReview(c context.Context, id string) error
}

type reviewService struct {
	repository repository.ReviewRepository
}

func NewReviewService(r repository.ReviewRepository) ReviewService {
	return &reviewService{r}
}

func (s *reviewService) PutReview(c context.Context, novelID, accountID string, rating int, title, content string, isSpoiler bool) (*model.Review, error) {
	rev := &model.Review{
		ID:        ksuid.New().String(),
		NovelID:   novelID,
		AccountID: accountID,
		Rating:    rating,
		Title:     title,
		Content:   content,
		IsSpoiler: isSpoiler,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repository.PutReview(c, rev); err != nil {
		return nil, err
	}

	return rev, nil
}

func (s *reviewService) GetReviewById(c context.Context, id string) (*model.Review, error) {
	return s.repository.GetReviewById(c, id)
}

func (s *reviewService) GetReviewsByNovel(c context.Context, novelID string, skip, take uint64) ([]*model.Review, error) {
	return s.repository.GetReviewsByNovel(c, novelID, skip, take)
}

func (s *reviewService) DeleteReview(c context.Context, id string) error {
	return s.repository.DeleteReview(c, id)
}