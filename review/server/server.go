package server

import (
	"context"
	"fmt"
	"net"

	"github.com/wignn/micro-3/review/genproto"
	"github.com/wignn/micro-3/review/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service service.ReviewService
	genproto.UnimplementedReviewServiceServer
}

func ListenGRPC(s service.ReviewService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	genproto.RegisterReviewServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostReview(ctx context.Context, req *genproto.PostReviewRequest) (*genproto.ReviewResponse, error) {
	rev, err := s.service.PutReview(ctx, req.NovelId, req.AccountId, int(req.Rating), req.Title, req.Content, req.IsSpoiler)
	if err != nil {
		return nil, err
	}
	createdAtBytes, _ := rev.CreatedAt.MarshalBinary()
	return &genproto.ReviewResponse{Review: &genproto.Review{
		Id:        rev.ID,
		NovelId:   rev.NovelID,
		AccountId: rev.AccountID,
		Rating:    int32(rev.Rating),
		Title:     rev.Title,
		Content:   rev.Content,
		IsSpoiler: rev.IsSpoiler,
		CreatedAt: createdAtBytes,
	}}, nil
}

func (s *grpcServer) GetReview(ctx context.Context, req *genproto.ReviewIdRequest) (*genproto.Review, error) {
	rev, err := s.service.GetReviewById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if rev == nil {
		return nil, fmt.Errorf("review not found")
	}
	createdAtBytes, _ := rev.CreatedAt.MarshalBinary()
	return &genproto.Review{
		Id:        rev.ID,
		NovelId:   rev.NovelID,
		AccountId: rev.AccountID,
		Rating:    int32(rev.Rating),
		Title:     rev.Title,
		Content:   rev.Content,
		IsSpoiler: rev.IsSpoiler,
		Upvotes:   rev.Upvotes,
		Downvotes: rev.Downvotes,
		CreatedAt: createdAtBytes,
	}, nil
}

func (s *grpcServer) GetReviewsByNovel(ctx context.Context, req *genproto.NovelReviewsRequest) (*genproto.ReviewListResponse, error) {
	reviews, err := s.service.GetReviewsByNovel(ctx, req.NovelId, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var protoReviews []*genproto.Review
	for _, rev := range reviews {
		createdAtBytes, _ := rev.CreatedAt.MarshalBinary()
		protoReviews = append(protoReviews, &genproto.Review{
			Id:        rev.ID,
			NovelId:   rev.NovelID,
			AccountId: rev.AccountID,
			Rating:    int32(rev.Rating),
			Title:     rev.Title,
			Content:   rev.Content,
			IsSpoiler: rev.IsSpoiler,
			Upvotes:   rev.Upvotes,
			Downvotes: rev.Downvotes,
			CreatedAt: createdAtBytes,
		})
	}
	return &genproto.ReviewListResponse{Reviews: protoReviews}, nil
}

func (s *grpcServer) DeleteReview(ctx context.Context, req *genproto.DeleteReviewRequest) (*genproto.DeleteReviewResponse, error) {
	if err := s.service.DeleteReview(ctx, req.Id); err != nil {
		return nil, err
	}
	return &genproto.DeleteReviewResponse{
		DeletedId: req.Id, Message: "Review deleted", Success: true,
	}, nil
}
