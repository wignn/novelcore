package client

import (
	"context"
	"log"
	"time"

	"github.com/wignn/micro-3/novel/genproto"
	"github.com/wignn/micro-3/novel/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NovelClient struct {
	conn    *grpc.ClientConn
	service genproto.NovelServiceClient
}

func NewClient(url string) (*NovelClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := genproto.NewNovelServiceClient(conn)
	return &NovelClient{conn, c}, nil
}

func (cl *NovelClient) Close() {
	cl.conn.Close()
}

// ── Novel ──────────────────────────────

func (cl *NovelClient) CreateNovel(c context.Context, req *genproto.CreateNovelRequest) (*genproto.Novel, error) {
	r, err := cl.service.CreateNovel(c, req)
	if err != nil {
		log.Printf("failed to create novel: %v\n", err)
		return nil, err
	}
	return r.Novel, nil
}

func (cl *NovelClient) GetNovel(c context.Context, id string) (*genproto.Novel, error) {
	r, err := cl.service.GetNovel(c, &genproto.GetNovelRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return r.Novel, nil
}

func (cl *NovelClient) ListNovels(c context.Context, req *genproto.ListNovelsRequest) ([]*genproto.Novel, error) {
	r, err := cl.service.ListNovels(c, req)
	if err != nil {
		return nil, err
	}
	return r.Novels, nil
}

func (cl *NovelClient) UpdateNovel(c context.Context, req *genproto.UpdateNovelRequest) (*genproto.Novel, error) {
	r, err := cl.service.UpdateNovel(c, req)
	if err != nil {
		return nil, err
	}
	return r.Novel, nil
}

func (cl *NovelClient) DeleteNovel(c context.Context, id string) (*genproto.DeleteResponse, error) {
	return cl.service.DeleteNovel(c, &genproto.DeleteRequest{Id: id})
}

// ── Chapter ────────────────────────────

func (cl *NovelClient) CreateChapter(c context.Context, req *genproto.CreateChapterRequest) (*genproto.Chapter, error) {
	r, err := cl.service.CreateChapter(c, req)
	if err != nil {
		return nil, err
	}
	return r.Chapter, nil
}

func (cl *NovelClient) GetChapter(c context.Context, id string) (*genproto.Chapter, error) {
	r, err := cl.service.GetChapter(c, &genproto.GetChapterRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return r.Chapter, nil
}

func (cl *NovelClient) ListChapters(c context.Context, novelID string, skip, take uint64) ([]*genproto.Chapter, error) {
	r, err := cl.service.ListChapters(c, &genproto.ListChaptersRequest{
		NovelId: novelID, Skip: skip, Take: take,
	})
	if err != nil {
		return nil, err
	}
	return r.Chapters, nil
}

func (cl *NovelClient) UpdateChapter(c context.Context, req *genproto.UpdateChapterRequest) (*genproto.Chapter, error) {
	r, err := cl.service.UpdateChapter(c, req)
	if err != nil {
		return nil, err
	}
	return r.Chapter, nil
}

func (cl *NovelClient) DeleteChapter(c context.Context, id string) (*genproto.DeleteResponse, error) {
	return cl.service.DeleteChapter(c, &genproto.DeleteRequest{Id: id})
}

// ── Author ─────────────────────────────

func (cl *NovelClient) CreateAuthor(c context.Context, name, bio string) (*genproto.Author, error) {
	r, err := cl.service.CreateAuthor(c, &genproto.CreateAuthorRequest{Name: name, Bio: bio})
	if err != nil {
		return nil, err
	}
	return r.Author, nil
}

func (cl *NovelClient) ListAuthors(c context.Context, skip, take uint64, id string) ([]*genproto.Author, error) {
	r, err := cl.service.ListAuthors(c, &genproto.ListAuthorsRequest{Skip: skip, Take: take, Id: id})
	if err != nil {
		return nil, err
	}
	return r.Authors, nil
}

// ── Translation Group ──────────────────

func (cl *NovelClient) CreateTranslationGroup(c context.Context, name, websiteURL, desc string) (*genproto.TranslationGroup, error) {
	r, err := cl.service.CreateTranslationGroup(c, &genproto.CreateTranslationGroupRequest{
		Name: name, WebsiteUrl: websiteURL, Description: desc,
	})
	if err != nil {
		return nil, err
	}
	return r.Group, nil
}

func (cl *NovelClient) ListTranslationGroups(c context.Context, skip, take uint64) ([]*genproto.TranslationGroup, error) {
	r, err := cl.service.ListTranslationGroups(c, &genproto.ListTranslationGroupsRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}
	return r.Groups, nil
}

// ── Genre & Tag ────────────────────────

func (cl *NovelClient) GetGenres(c context.Context) ([]*genproto.Genre, error) {
	r, err := cl.service.GetGenres(c, &genproto.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return r.Genres, nil
}

func (cl *NovelClient) GetTags(c context.Context) ([]*genproto.Tag, error) {
	r, err := cl.service.GetTags(c, &genproto.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return r.Tags, nil
}

// ── Ranking ────────────────────────────

func (cl *NovelClient) GetRanking(c context.Context, period, sortBy string, skip, take uint64) ([]*genproto.RankedNovel, error) {
	r, err := cl.service.GetRanking(c, &genproto.RankingRequest{
		Period: period, SortBy: sortBy, Skip: skip, Take: take,
	})
	if err != nil {
		return nil, err
	}
	return r.RankedNovels, nil
}

// ── View ───────────────────────────────

func (cl *NovelClient) IncrementViewCount(c context.Context, novelID string) (int64, error) {
	r, err := cl.service.IncrementViewCount(c, &genproto.IncrementViewRequest{NovelId: novelID})
	if err != nil {
		return 0, err
	}
	return r.ViewCount, nil
}

// ── Helpers ────────────────────────────

func ProtoToNovelModel(n *genproto.Novel) *model.Novel {
	if n == nil {
		return nil
	}

	var createdAt, updatedAt time.Time
	if n.CreatedAt != nil {
		createdAt.UnmarshalBinary(n.CreatedAt)
	}
	if n.UpdatedAt != nil {
		updatedAt.UnmarshalBinary(n.UpdatedAt)
	}

	novel := &model.Novel{
		ID:               n.Id,
		Title:            n.Title,
		AlternativeTitle: n.AlternativeTitle,
		Description:      n.Description,
		CoverImageURL:    n.CoverImageUrl,
		AuthorID:         n.AuthorId,
		Status:           n.Status,
		NovelType:        n.NovelType,
		CountryOfOrigin:  n.CountryOfOrigin,
		YearPublished:    n.YearPublished,
		TotalChapters:    n.TotalChapters,
		RatingAvg:        n.RatingAvg,
		RatingCount:      n.RatingCount,
		ViewCount:        n.ViewCount,
		BookmarkCount:    n.BookmarkCount,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}

	for _, g := range n.Genres {
		novel.Genres = append(novel.Genres, model.Genre{ID: g.Id, Name: g.Name, Slug: g.Slug})
	}
	for _, t := range n.Tags {
		novel.Tags = append(novel.Tags, model.Tag{ID: t.Id, Name: t.Name, Slug: t.Slug})
	}

	if n.Author != nil {
		var authorCreatedAt time.Time
		if n.Author.CreatedAt != nil {
			authorCreatedAt.UnmarshalBinary(n.Author.CreatedAt)
		}
		novel.Author = &model.Author{
			ID: n.Author.Id, Name: n.Author.Name, Bio: n.Author.Bio, CreatedAt: authorCreatedAt,
		}
	}

	return novel
}
