package client

import (
	"context"
	"time"

	"github.com/wignn/micro-3/review/genproto"
	"github.com/wignn/micro-3/review/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReviewClient struct {
	conn    *grpc.ClientConn
	service genproto.ReviewServiceClient
}

func NewClient(url string) (*ReviewClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := genproto.NewReviewServiceClient(conn)
	return &ReviewClient{conn, c}, nil
}

func (cl *ReviewClient) Close() {
	cl.conn.Close()
}

func (cl *ReviewClient) PostReview(c context.Context, novelID, accountID, title, content string, rating int32, isSpoiler bool) (*model.Review, error) {
	r, err := cl.service.PostReview(c, &genproto.PostReviewRequest{
		NovelId:   novelID,
		AccountId: accountID,
		Title:     title,
		Content:   content,
		Rating:    rating,
		IsSpoiler: isSpoiler,
	})
	if err != nil {
		return nil, err
	}
	var createdAt time.Time
	if r.Review.CreatedAt != nil {
		createdAt.UnmarshalBinary(r.Review.CreatedAt)
	}
	return &model.Review{
		ID:        r.Review.Id,
		NovelID:   r.Review.NovelId,
		AccountID: r.Review.AccountId,
		Rating:    int(r.Review.Rating),
		Title:     r.Review.Title,
		Content:   r.Review.Content,
		IsSpoiler: r.Review.IsSpoiler,
		CreatedAt: createdAt,
	}, nil
}

func (cl *ReviewClient) GetReview(c context.Context, id string) (*genproto.Review, error) {
	return cl.service.GetReview(c, &genproto.ReviewIdRequest{Id: id})
}

func (cl *ReviewClient) GetReviewsByNovel(c context.Context, novelID string, skip, take uint64) ([]*genproto.Review, error) {
	r, err := cl.service.GetReviewsByNovel(c, &genproto.NovelReviewsRequest{
		NovelId: novelID, Skip: skip, Take: take,
	})
	if err != nil {
		return nil, err
	}
	return r.Reviews, nil
}

func (cl *ReviewClient) DeleteReview(c context.Context, id string) (*genproto.DeleteReviewResponse, error) {
	return cl.service.DeleteReview(c, &genproto.DeleteReviewRequest{Id: id})
}
