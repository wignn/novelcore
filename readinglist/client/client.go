package client

import (
	"context"
	"time"

	"github.com/wignn/micro-3/readinglist/genproto"
	"github.com/wignn/micro-3/readinglist/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReadingListClient struct {
	conn    *grpc.ClientConn
	service genproto.ReadingListServiceClient
}

func NewClient(url string) (*ReadingListClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := genproto.NewReadingListServiceClient(conn)
	return &ReadingListClient{conn, c}, nil
}

func (cl *ReadingListClient) Close() {
	cl.conn.Close()
}

func (cl *ReadingListClient) AddToReadingList(c context.Context, accountID, novelID, status string, currentChapter float64, rating int32, notes string, isFavorite bool) (*model.ReadingListEntry, error) {
	r, err := cl.service.AddToReadingList(c, &genproto.AddToReadingListRequest{
		AccountId:      accountID,
		NovelId:        novelID,
		Status:         status,
		CurrentChapter: currentChapter,
		Rating:         rating,
		Notes:          notes,
		IsFavorite:     isFavorite,
	})
	if err != nil {
		return nil, err
	}
	return protoToEntry(r.Entry), nil
}

func (cl *ReadingListClient) UpdateReadingList(c context.Context, id, status string, currentChapter float64, rating int32, notes string, isFavorite bool) (*model.ReadingListEntry, error) {
	r, err := cl.service.UpdateReadingList(c, &genproto.UpdateReadingListRequest{
		Id:             id,
		Status:         status,
		CurrentChapter: currentChapter,
		Rating:         rating,
		Notes:          notes,
		IsFavorite:     isFavorite,
	})
	if err != nil {
		return nil, err
	}
	return protoToEntry(r.Entry), nil
}

func (cl *ReadingListClient) GetReadingList(c context.Context, accountID, status string, skip, take uint64) ([]*model.ReadingListEntry, error) {
	r, err := cl.service.GetReadingList(c, &genproto.GetReadingListRequest{
		AccountId: accountID, Status: status, Skip: skip, Take: take,
	})
	if err != nil {
		return nil, err
	}
	var entries []*model.ReadingListEntry
	for _, e := range r.Entries {
		entries = append(entries, protoToEntry(e))
	}
	return entries, nil
}

func (cl *ReadingListClient) RemoveFromReadingList(c context.Context, id string) (*genproto.DeleteReadingListResponse, error) {
	return cl.service.RemoveFromReadingList(c, &genproto.DeleteReadingListRequest{Id: id})
}

func protoToEntry(e *genproto.ReadingListEntry) *model.ReadingListEntry {
	var createdAt, updatedAt time.Time
	if e.CreatedAt != nil {
		createdAt.UnmarshalBinary(e.CreatedAt)
	}
	if e.UpdatedAt != nil {
		updatedAt.UnmarshalBinary(e.UpdatedAt)
	}
	return &model.ReadingListEntry{
		ID:             e.Id,
		AccountID:      e.AccountId,
		NovelID:        e.NovelId,
		Status:         e.Status,
		CurrentChapter: e.CurrentChapter,
		Rating:         e.Rating,
		Notes:          e.Notes,
		IsFavorite:     e.IsFavorite,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
