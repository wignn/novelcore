package server

import (
	"context"
	"fmt"
	"net"

	"github.com/wignn/micro-3/readinglist/genproto"
	"github.com/wignn/micro-3/readinglist/model"
	"github.com/wignn/micro-3/readinglist/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service service.ReadingListService
	genproto.UnimplementedReadingListServiceServer
}

func ListenGRPC(s service.ReadingListService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	genproto.RegisterReadingListServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) AddToReadingList(c context.Context, req *genproto.AddToReadingListRequest) (*genproto.ReadingListResponse, error) {
	e, err := s.service.AddEntry(c, req.AccountId, req.NovelId, req.Status, req.CurrentChapter, req.Rating, req.Notes, req.IsFavorite)
	if err != nil {
		return nil, err
	}
	return &genproto.ReadingListResponse{Entry: entryToProto(e)}, nil
}

func (s *grpcServer) UpdateReadingList(c context.Context, req *genproto.UpdateReadingListRequest) (*genproto.ReadingListResponse, error) {
	e, err := s.service.UpdateEntry(c, req.Id, req.Status, req.CurrentChapter, req.Rating, req.Notes, req.IsFavorite)
	if err != nil {
		return nil, err
	}
	return &genproto.ReadingListResponse{Entry: entryToProto(e)}, nil
}

func (s *grpcServer) GetReadingList(c context.Context, req *genproto.GetReadingListRequest) (*genproto.GetReadingListResponse, error) {
	entries, err := s.service.GetEntries(c, req.AccountId, req.Status, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var protoEntries []*genproto.ReadingListEntry
	for _, e := range entries {
		protoEntries = append(protoEntries, entryToProto(e))
	}
	return &genproto.GetReadingListResponse{Entries: protoEntries}, nil
}

func (s *grpcServer) RemoveFromReadingList(c context.Context, req *genproto.DeleteReadingListRequest) (*genproto.DeleteReadingListResponse, error) {
	if err := s.service.RemoveEntry(c, req.Id); err != nil {
		return nil, err
	}
	return &genproto.DeleteReadingListResponse{
		DeletedId: req.Id, Message: "Removed from reading list", Success: true,
	}, nil
}

func entryToProto(e *model.ReadingListEntry) *genproto.ReadingListEntry {
	createdAt, _ := e.CreatedAt.MarshalBinary()
	updatedAt, _ := e.UpdatedAt.MarshalBinary()
	return &genproto.ReadingListEntry{
		Id:             e.ID,
		AccountId:      e.AccountID,
		NovelId:        e.NovelID,
		Status:         e.Status,
		CurrentChapter: e.CurrentChapter,
		Rating:         e.Rating,
		Notes:          e.Notes,
		IsFavorite:     e.IsFavorite,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
