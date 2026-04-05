package main

import (
	"context"
	"log"
	"strings"
	"time"

	novelGenproto "github.com/wignn/micro-3/novel/genproto"
)

type queryResolver struct {
	server *GraphQLServer
}

// ── Account ────────────────────────────

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		a, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Account{{
			ID: a.ID, Name: a.Name, Email: a.Email,
			AvatarUrl: &a.AvatarURL, Bio: &a.Bio, Role: a.Role, CreatedAt: a.CreatedAt,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	for _, a := range accountList {
		accounts = append(accounts, &Account{
			ID: a.ID, Name: a.Name, Email: a.Email,
			AvatarUrl: &a.AvatarURL, Bio: &a.Bio, Role: a.Role, CreatedAt: a.CreatedAt,
		})
	}
	return accounts, nil
}

// ── Novel ──────────────────────────────

func (r *queryResolver) Novels(ctx context.Context, pagination *PaginationInput, id *string, filter *NovelFilterInput, query *string) ([]*Novel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if id != nil {
		n, err := r.server.novelClient.GetNovel(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Novel{protoNovelToGraphQL(n)}, nil
	}

	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	req := &novelGenproto.ListNovelsRequest{Skip: skip, Take: take}

	if query != nil {
		req.Query = *query
	}
	if filter != nil {
		if filter.Status != nil {
			req.Status = *filter.Status
		}
		if filter.NovelType != nil {
			req.NovelType = *filter.NovelType
		}
		if filter.CountryOfOrigin != nil {
			req.CountryOfOrigin = *filter.CountryOfOrigin
		}
		if filter.SortBy != nil {
			req.SortBy = *filter.SortBy
		}
		if filter.SortOrder != nil {
			req.SortOrder = *filter.SortOrder
		}
		for _, gid := range filter.GenreIds {
			req.GenreIds = append(req.GenreIds, int32(gid))
		}
		for _, tid := range filter.TagIds {
			req.TagIds = append(req.TagIds, int32(tid))
		}
	}

	novels, err := r.server.novelClient.ListNovels(ctx, req)
	if err != nil {
		return nil, err
	}

	var result []*Novel
	for _, n := range novels {
		result = append(result, protoNovelToGraphQL(n))
	}
	return result, nil
}

// ── Chapter ────────────────────────────

func (r *queryResolver) Chapters(ctx context.Context, novelID string, pagination *PaginationInput) ([]*Chapter, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	skip, take := uint64(0), uint64(500)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	chapters, err := r.server.novelClient.ListChapters(ctx, novelID, skip, take)
	if err != nil {
		return nil, err
	}

	var result []*Chapter
	for _, ch := range chapters {
		result = append(result, protoChapterToGraphQL(ch))
	}
	return result, nil
}

func (r *queryResolver) Chapter(ctx context.Context, id string) (*Chapter, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ch, err := r.server.novelClient.GetChapter(ctx, id)
	if err != nil {
		return nil, err
	}
	return protoChapterToGraphQL(ch), nil
}

// ── Author ─────────────────────────────

func (r *queryResolver) Authors(ctx context.Context, pagination *PaginationInput, id *string) ([]*Author, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	authorID := ""
	if id != nil {
		authorID = *id
	}

	authors, err := r.server.novelClient.ListAuthors(ctx, skip, take, authorID)
	if err != nil {
		return nil, err
	}

	var result []*Author
	for _, a := range authors {
		var createdAt time.Time
		if a.CreatedAt != nil {
			createdAt.UnmarshalBinary(a.CreatedAt)
		}
		result = append(result, &Author{ID: a.Id, Name: a.Name, Bio: &a.Bio, CreatedAt: createdAt})
	}
	return result, nil
}

// ── Translation Group ──────────────────

func (r *queryResolver) TranslationGroups(ctx context.Context, pagination *PaginationInput) ([]*TranslationGroup, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	groups, err := r.server.novelClient.ListTranslationGroups(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var result []*TranslationGroup
	for _, g := range groups {
		var createdAt time.Time
		if g.CreatedAt != nil {
			createdAt.UnmarshalBinary(g.CreatedAt)
		}
		result = append(result, &TranslationGroup{
			ID: g.Id, Name: g.Name, WebsiteURL: &g.WebsiteUrl,
			Description: &g.Description, CreatedAt: createdAt,
		})
	}
	return result, nil
}

// ── Genre & Tag ────────────────────────

func (r *queryResolver) Genres(ctx context.Context) ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	genres, err := r.server.novelClient.GetGenres(ctx)
	if err != nil {
		return nil, err
	}

	var result []*Genre
	for _, g := range genres {
		result = append(result, &Genre{ID: int(g.Id), Name: g.Name, Slug: g.Slug})
	}
	return result, nil
}

func (r *queryResolver) Tags(ctx context.Context) ([]*Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tags, err := r.server.novelClient.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	var result []*Tag
	for _, t := range tags {
		result = append(result, &Tag{ID: int(t.Id), Name: t.Name, Slug: t.Slug})
	}
	return result, nil
}

// ── Reading List ───────────────────────

func (r *queryResolver) ReadingList(ctx context.Context, accountID string, status *string, pagination *PaginationInput) ([]*ReadingListEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	statusFilter := ""
	if status != nil {
		statusFilter = *status
	}

	entries, err := r.server.readinglistClient.GetReadingList(ctx, accountID, statusFilter, skip, take)
	if err != nil {
		return nil, err
	}

	var result []*ReadingListEntry
	for _, e := range entries {
		ratingInt := int(e.Rating)
		result = append(result, &ReadingListEntry{
			ID: e.ID, NovelID: e.NovelID, Status: e.Status,
			CurrentChapter: e.CurrentChapter, Rating: &ratingInt,
			Notes: &e.Notes, IsFavorite: e.IsFavorite,
			CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
		})
	}
	return result, nil
}

// ── Review ─────────────────────────────

func (r *queryResolver) Reviews(ctx context.Context, novelID string, pagination *PaginationInput) ([]*Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	reviews, err := r.server.reviewClient.GetReviewsByNovel(ctx, novelID, skip, take)
	if err != nil {
		return nil, err
	}

	var result []*Review
	for _, rv := range reviews {
		var createdAt time.Time
		if rv.CreatedAt != nil {
			createdAt.UnmarshalBinary(rv.CreatedAt)
		}

		result = append(result, &Review{
			ID: rv.Id, NovelID: rv.NovelId, AccountID: rv.AccountId,
			Rating: int(rv.Rating), Title: &rv.Title, Content: &rv.Content,
			IsSpoiler: rv.IsSpoiler, Upvotes: int(rv.Upvotes), Downvotes: int(rv.Downvotes),
			CreatedAt: createdAt,
		})
	}
	return result, nil
}

// ── Ranking ────────────────────────────

func (r *queryResolver) NovelRanking(ctx context.Context, period RankingPeriod, sortBy RankingSortBy, pagination *PaginationInput) ([]*RankedNovel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	skip, take := uint64(0), uint64(50)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	periodStr := strings.ToLower(string(period))
	sortByStr := strings.ToLower(string(sortBy))

	ranked, err := r.server.novelClient.GetRanking(ctx, periodStr, sortByStr, skip, take)
	if err != nil {
		log.Println("GetRanking error:", err)
		return nil, err
	}

	var result []*RankedNovel
	for _, rn := range ranked {
		result = append(result, &RankedNovel{
			Rank:  int(rn.Rank),
			Novel: protoNovelToGraphQL(rn.Novel),
			Score: rn.Score,
		})
	}
	return result, nil
}

// ── Helpers ────────────────────────────

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue
}
