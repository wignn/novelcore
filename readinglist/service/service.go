package service

import (
	"context"

	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/readinglist/model"
	"github.com/wignn/micro-3/readinglist/repository"
)

type ReadingListService interface {
	AddEntry(c context.Context, accountID, novelID, status string, currentChapter float64, rating int32, notes string, isFavorite bool) (*model.ReadingListEntry, error)
	UpdateEntry(c context.Context, id, status string, currentChapter float64, rating int32, notes string, isFavorite bool) (*model.ReadingListEntry, error)
	GetEntries(c context.Context, accountID, status string, skip, take uint64) ([]*model.ReadingListEntry, error)
	RemoveEntry(c context.Context, id string) error
}

type readingListService struct {
	repo repository.ReadingListRepository
}

func NewReadingListService(r repository.ReadingListRepository) ReadingListService {
	return &readingListService{repo: r}
}

func (s *readingListService) AddEntry(c context.Context, accountID, novelID, status string, currentChapter float64, rating int32, notes string, isFavorite bool) (*model.ReadingListEntry, error) {
	if status == "" {
		status = "plan_to_read"
	}

	e := &model.ReadingListEntry{
		ID:             ksuid.New().String(),
		AccountID:      accountID,
		NovelID:        novelID,
		Status:         status,
		CurrentChapter: currentChapter,
		Rating:         rating,
		Notes:          notes,
		IsFavorite:     isFavorite,
	}

	if err := s.repo.AddEntry(c, e); err != nil {
		return nil, err
	}

	return s.repo.GetEntryByID(c, e.ID)
}

func (s *readingListService) UpdateEntry(c context.Context, id, status string, currentChapter float64, rating int32, notes string, isFavorite bool) (*model.ReadingListEntry, error) {
	e := &model.ReadingListEntry{
		ID:             id,
		Status:         status,
		CurrentChapter: currentChapter,
		Rating:         rating,
		Notes:          notes,
		IsFavorite:     isFavorite,
	}

	if err := s.repo.UpdateEntry(c, e); err != nil {
		return nil, err
	}

	return s.repo.GetEntryByID(c, id)
}

func (s *readingListService) GetEntries(c context.Context, accountID, status string, skip, take uint64) ([]*model.ReadingListEntry, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}
	return s.repo.GetEntries(c, accountID, status, skip, take)
}

func (s *readingListService) RemoveEntry(c context.Context, id string) error {
	return s.repo.DeleteEntry(c, id)
}
