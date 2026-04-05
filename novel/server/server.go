package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/wignn/micro-3/novel/genproto"
	"github.com/wignn/micro-3/novel/model"
	"github.com/wignn/micro-3/novel/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service service.NovelService
	genproto.UnimplementedNovelServiceServer
}

func ListenGRPC(s service.NovelService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	genproto.RegisterNovelServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

// ── Novel ──────────────────────────────

func (s *grpcServer) CreateNovel(c context.Context, req *genproto.CreateNovelRequest) (*genproto.NovelResponse, error) {
	n, err := s.service.CreateNovel(c, req.Title, req.AlternativeTitle, req.Description,
		req.CoverImageUrl, req.AuthorId, req.Status, req.NovelType, req.CountryOfOrigin,
		req.YearPublished, req.GenreIds, req.TagIds)
	if err != nil {
		return nil, err
	}
	return &genproto.NovelResponse{Novel: novelToProto(n)}, nil
}

func (s *grpcServer) GetNovel(c context.Context, req *genproto.GetNovelRequest) (*genproto.NovelResponse, error) {
	n, err := s.service.GetNovel(c, req.Id)
	if err != nil {
		return nil, err
	}
	return &genproto.NovelResponse{Novel: novelToProto(n)}, nil
}

func (s *grpcServer) ListNovels(c context.Context, req *genproto.ListNovelsRequest) (*genproto.ListNovelsResponse, error) {
	novels, err := s.service.ListNovels(c, req.Skip, req.Take, req.Query, req.Status,
		req.NovelType, req.CountryOfOrigin, req.SortBy, req.SortOrder, req.GenreIds, req.TagIds)
	if err != nil {
		return nil, err
	}

	var protoNovels []*genproto.Novel
	for _, n := range novels {
		protoNovels = append(protoNovels, novelToProto(n))
	}
	return &genproto.ListNovelsResponse{Novels: protoNovels}, nil
}

func (s *grpcServer) UpdateNovel(c context.Context, req *genproto.UpdateNovelRequest) (*genproto.NovelResponse, error) {
	n, err := s.service.UpdateNovel(c, req.Id, req.Title, req.AlternativeTitle, req.Description,
		req.CoverImageUrl, req.AuthorId, req.Status, req.NovelType, req.CountryOfOrigin,
		req.YearPublished, req.GenreIds, req.TagIds)
	if err != nil {
		return nil, err
	}
	return &genproto.NovelResponse{Novel: novelToProto(n)}, nil
}

func (s *grpcServer) DeleteNovel(c context.Context, req *genproto.DeleteRequest) (*genproto.DeleteResponse, error) {
	if err := s.service.DeleteNovel(c, req.Id); err != nil {
		return nil, err
	}
	return &genproto.DeleteResponse{DeletedId: req.Id, Message: "Novel deleted", Success: true}, nil
}

// ── Chapter ────────────────────────────

func (s *grpcServer) CreateChapter(c context.Context, req *genproto.CreateChapterRequest) (*genproto.ChapterResponse, error) {
	ch, err := s.service.CreateChapter(c, req.NovelId, req.ChapterNumber, req.Title, req.TranslatorGroupId, req.SourceUrl)
	if err != nil {
		return nil, err
	}
	return &genproto.ChapterResponse{Chapter: chapterToProto(ch)}, nil
}

func (s *grpcServer) GetChapter(c context.Context, req *genproto.GetChapterRequest) (*genproto.ChapterResponse, error) {
	ch, err := s.service.GetChapter(c, req.Id)
	if err != nil {
		return nil, err
	}
	return &genproto.ChapterResponse{Chapter: chapterToProto(ch)}, nil
}

func (s *grpcServer) ListChapters(c context.Context, req *genproto.ListChaptersRequest) (*genproto.ListChaptersResponse, error) {
	chapters, err := s.service.ListChapters(c, req.NovelId, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var protoChapters []*genproto.Chapter
	for _, ch := range chapters {
		protoChapters = append(protoChapters, chapterToProto(ch))
	}
	return &genproto.ListChaptersResponse{Chapters: protoChapters}, nil
}

func (s *grpcServer) UpdateChapter(c context.Context, req *genproto.UpdateChapterRequest) (*genproto.ChapterResponse, error) {
	ch, err := s.service.UpdateChapter(c, req.Id, req.NovelId, req.ChapterNumber, req.Title, req.TranslatorGroupId, req.SourceUrl)
	if err != nil {
		return nil, err
	}
	return &genproto.ChapterResponse{Chapter: chapterToProto(ch)}, nil
}

func (s *grpcServer) DeleteChapter(c context.Context, req *genproto.DeleteRequest) (*genproto.DeleteResponse, error) {
	if err := s.service.DeleteChapter(c, req.Id); err != nil {
		return nil, err
	}
	return &genproto.DeleteResponse{DeletedId: req.Id, Message: "Chapter deleted", Success: true}, nil
}

// ── Author ─────────────────────────────

func (s *grpcServer) CreateAuthor(c context.Context, req *genproto.CreateAuthorRequest) (*genproto.AuthorResponse, error) {
	a, err := s.service.CreateAuthor(c, req.Name, req.Bio)
	if err != nil {
		return nil, err
	}
	createdAt, _ := a.CreatedAt.MarshalBinary()
	return &genproto.AuthorResponse{
		Author: &genproto.Author{Id: a.ID, Name: a.Name, Bio: a.Bio, CreatedAt: createdAt},
	}, nil
}

func (s *grpcServer) ListAuthors(c context.Context, req *genproto.ListAuthorsRequest) (*genproto.ListAuthorsResponse, error) {
	authors, err := s.service.ListAuthors(c, req.Skip, req.Take, req.Id)
	if err != nil {
		return nil, err
	}
	var protoAuthors []*genproto.Author
	for _, a := range authors {
		createdAt, _ := a.CreatedAt.MarshalBinary()
		protoAuthors = append(protoAuthors, &genproto.Author{
			Id: a.ID, Name: a.Name, Bio: a.Bio, CreatedAt: createdAt,
		})
	}
	return &genproto.ListAuthorsResponse{Authors: protoAuthors}, nil
}

// ── Translation Group ──────────────────

func (s *grpcServer) CreateTranslationGroup(c context.Context, req *genproto.CreateTranslationGroupRequest) (*genproto.TranslationGroupResponse, error) {
	g, err := s.service.CreateTranslationGroup(c, req.Name, req.WebsiteUrl, req.Description)
	if err != nil {
		return nil, err
	}
	createdAt, _ := g.CreatedAt.MarshalBinary()
	return &genproto.TranslationGroupResponse{
		Group: &genproto.TranslationGroup{
			Id: g.ID, Name: g.Name, WebsiteUrl: g.WebsiteURL, Description: g.Description, CreatedAt: createdAt,
		},
	}, nil
}

func (s *grpcServer) ListTranslationGroups(c context.Context, req *genproto.ListTranslationGroupsRequest) (*genproto.ListTranslationGroupsResponse, error) {
	groups, err := s.service.ListTranslationGroups(c, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var protoGroups []*genproto.TranslationGroup
	for _, g := range groups {
		createdAt, _ := g.CreatedAt.MarshalBinary()
		protoGroups = append(protoGroups, &genproto.TranslationGroup{
			Id: g.ID, Name: g.Name, WebsiteUrl: g.WebsiteURL, Description: g.Description, CreatedAt: createdAt,
		})
	}
	return &genproto.ListTranslationGroupsResponse{Groups: protoGroups}, nil
}

// ── Genre & Tag ────────────────────────

func (s *grpcServer) GetGenres(c context.Context, _ *genproto.EmptyRequest) (*genproto.GenreListResponse, error) {
	genres, err := s.service.GetGenres(c)
	if err != nil {
		return nil, err
	}
	var protoGenres []*genproto.Genre
	for _, g := range genres {
		protoGenres = append(protoGenres, &genproto.Genre{Id: g.ID, Name: g.Name, Slug: g.Slug})
	}
	return &genproto.GenreListResponse{Genres: protoGenres}, nil
}

func (s *grpcServer) GetTags(c context.Context, _ *genproto.EmptyRequest) (*genproto.TagListResponse, error) {
	tags, err := s.service.GetTags(c)
	if err != nil {
		return nil, err
	}
	var protoTags []*genproto.Tag
	for _, t := range tags {
		protoTags = append(protoTags, &genproto.Tag{Id: t.ID, Name: t.Name, Slug: t.Slug})
	}
	return &genproto.TagListResponse{Tags: protoTags}, nil
}

// ── Ranking ────────────────────────────

func (s *grpcServer) GetRanking(c context.Context, req *genproto.RankingRequest) (*genproto.RankingResponse, error) {
	novels, err := s.service.GetRanking(c, req.Period, req.SortBy, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var ranked []*genproto.RankedNovel
	for i, n := range novels {
		ranked = append(ranked, &genproto.RankedNovel{
			Rank:  int32(req.Skip) + int32(i) + 1,
			Novel: novelToProto(n),
			Score: n.RatingAvg,
		})
	}
	return &genproto.RankingResponse{RankedNovels: ranked}, nil
}

// ── View ───────────────────────────────

func (s *grpcServer) IncrementViewCount(c context.Context, req *genproto.IncrementViewRequest) (*genproto.IncrementViewResponse, error) {
	count, err := s.service.IncrementViewCount(c, req.NovelId)
	if err != nil {
		return nil, err
	}
	return &genproto.IncrementViewResponse{ViewCount: count}, nil
}

// ── Proto converters ───────────────────

func novelToProto(n *model.Novel) *genproto.Novel {
	if n == nil {
		return nil
	}

	createdAt, _ := n.CreatedAt.MarshalBinary()
	updatedAt, _ := n.UpdatedAt.MarshalBinary()

	p := &genproto.Novel{
		Id:              n.ID,
		Title:           n.Title,
		AlternativeTitle: n.AlternativeTitle,
		Description:     n.Description,
		CoverImageUrl:   n.CoverImageURL,
		AuthorId:        n.AuthorID,
		Status:          n.Status,
		NovelType:       n.NovelType,
		CountryOfOrigin: n.CountryOfOrigin,
		YearPublished:   n.YearPublished,
		TotalChapters:   n.TotalChapters,
		RatingAvg:       n.RatingAvg,
		RatingCount:     n.RatingCount,
		ViewCount:       n.ViewCount,
		BookmarkCount:   n.BookmarkCount,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	for _, g := range n.Genres {
		p.Genres = append(p.Genres, &genproto.Genre{Id: g.ID, Name: g.Name, Slug: g.Slug})
	}
	for _, t := range n.Tags {
		p.Tags = append(p.Tags, &genproto.Tag{Id: t.ID, Name: t.Name, Slug: t.Slug})
	}

	if n.Author != nil {
		authorCreatedAt, _ := n.Author.CreatedAt.MarshalBinary()
		p.Author = &genproto.Author{
			Id: n.Author.ID, Name: n.Author.Name, Bio: n.Author.Bio, CreatedAt: authorCreatedAt,
		}
	}

	log.Printf("Novel %s converted to proto", n.ID)

	return p
}

func chapterToProto(ch *model.Chapter) *genproto.Chapter {
	if ch == nil {
		return nil
	}
	createdAt, _ := ch.CreatedAt.MarshalBinary()
	updatedAt, _ := ch.UpdatedAt.MarshalBinary()

	p := &genproto.Chapter{
		Id:                ch.ID,
		NovelId:           ch.NovelID,
		ChapterNumber:     ch.ChapterNumber,
		Title:             ch.Title,
		TranslatorGroupId: ch.TranslatorGroupID,
		SourceUrl:         ch.SourceURL,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	if ch.TranslatorGroup != nil {
		tgCreatedAt, _ := ch.TranslatorGroup.CreatedAt.MarshalBinary()
		p.TranslatorGroup = &genproto.TranslationGroup{
			Id: ch.TranslatorGroup.ID, Name: ch.TranslatorGroup.Name,
			WebsiteUrl: ch.TranslatorGroup.WebsiteURL, Description: ch.TranslatorGroup.Description,
			CreatedAt: tgCreatedAt,
		}
	}

	return p
}
