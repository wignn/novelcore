package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	novelGenproto "github.com/wignn/micro-3/novel/genproto"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *GraphQLServer
}

func (r *mutationResolver) CreateAccount(c context.Context, in AccountInput) (*Account, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	a, err := r.server.accountClient.PostAccount(c, in.Name, in.Email, in.Password)
	if err != nil {
		return nil, handleError("CreateAccount", err)
	}

	return &Account{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		AvatarUrl: &a.AvatarURL,
		Bio:       &a.Bio,
		Role:      a.Role,
		CreatedAt: a.CreatedAt,
	}, nil
}

func (r *mutationResolver) Login(c context.Context, in LoginInput) (*AuthResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	token, err := r.server.authClient.Login(c, in.Email, in.Password)
	if err != nil {
		return nil, handleError("Login", err)
	}

	return &AuthResponse{
		ID:    token.Auth.Id,
		Email: token.Auth.Email,
		BackendToken: &Token{
			AccessToken:  token.Auth.Token.AccessToken,
			RefreshToken: token.Auth.Token.RefreshToken,
			ExpiresIn:    int(token.Auth.Token.ExpiresAt),
		},
	}, nil
}

func (r *mutationResolver) RefreshToken(c context.Context, refreshToken string) (*Token, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	newToken, err := r.server.authClient.RefreshToken(c, refreshToken)
	if err != nil {
		return nil, handleError("RefreshToken", err)
	}

	return &Token{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresIn:    int(newToken.ExpiresAt),
	}, nil
}

func (r *mutationResolver) EditAccount(c context.Context, id string, in EditAccountInput) (*Account, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	name, email, password, avatarUrl, bio := "", "", "", "", ""
	if in.Name != nil {
		name = *in.Name
	}
	if in.Email != nil {
		email = *in.Email
	}
	if in.Password != nil {
		password = *in.Password
	}
	if in.AvatarURL != nil {
		avatarUrl = *in.AvatarURL
	}
	if in.Bio != nil {
		bio = *in.Bio
	}

	a, err := r.server.accountClient.EditAccount(c, id, name, email, password, avatarUrl, bio)
	if err != nil {
		return nil, handleError("EditAccount", err)
	}

	return &Account{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		AvatarUrl: &a.AvatarURL,
		Bio:       &a.Bio,
		Role:      a.Role,
		CreatedAt: a.CreatedAt,
	}, nil
}

func (r *mutationResolver) DeleteAccount(c context.Context, id string) (*DeleteResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	p, err := r.server.accountClient.DeleteAccount(c, id)
	if err != nil {
		return nil, handleError("DeleteAccount", err)
	}

	return &DeleteResponse{
		DeletedID: p.DeletedID,
		Success:   p.Success,
		Message:   p.Message,
	}, nil
}


func (r *mutationResolver) CreateNovel(c context.Context, in NovelInput) (*Novel, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	req := &novelGenproto.CreateNovelRequest{
		Title: in.Title,
	}
	if in.AlternativeTitle != nil {
		req.AlternativeTitle = *in.AlternativeTitle
	}
	if in.Description != nil {
		req.Description = *in.Description
	}
	if in.CoverImageURL != nil {
		req.CoverImageUrl = *in.CoverImageURL
	}
	if in.AuthorID != nil {
		req.AuthorId = *in.AuthorID
	}
	if in.Status != nil {
		req.Status = *in.Status
	}
	if in.NovelType != nil {
		req.NovelType = *in.NovelType
	}
	if in.CountryOfOrigin != nil {
		req.CountryOfOrigin = *in.CountryOfOrigin
	}
	if in.YearPublished != nil {
		req.YearPublished = int32(*in.YearPublished)
	}

	for _, id := range in.GenreIds {
		req.GenreIds = append(req.GenreIds, int32(id))
	}
	for _, id := range in.TagIds {
		req.TagIds = append(req.TagIds, int32(id))
	}

	n, err := r.server.novelClient.CreateNovel(c, req)
	if err != nil {
		return nil, handleError("CreateNovel", err)
	}

	return protoNovelToGraphQL(n), nil
}

func (r *mutationResolver) UpdateNovel(c context.Context, id string, in NovelInput) (*Novel, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	req := &novelGenproto.UpdateNovelRequest{
		Id:    id,
		Title: in.Title,
	}
	if in.AlternativeTitle != nil {
		req.AlternativeTitle = *in.AlternativeTitle
	}
	if in.Description != nil {
		req.Description = *in.Description
	}
	if in.CoverImageURL != nil {
		req.CoverImageUrl = *in.CoverImageURL
	}
	if in.AuthorID != nil {
		req.AuthorId = *in.AuthorID
	}
	if in.Status != nil {
		req.Status = *in.Status
	}
	if in.NovelType != nil {
		req.NovelType = *in.NovelType
	}
	if in.CountryOfOrigin != nil {
		req.CountryOfOrigin = *in.CountryOfOrigin
	}
	if in.YearPublished != nil {
		req.YearPublished = int32(*in.YearPublished)
	}

	for _, gid := range in.GenreIds {
		req.GenreIds = append(req.GenreIds, int32(gid))
	}
	for _, tid := range in.TagIds {
		req.TagIds = append(req.TagIds, int32(tid))
	}

	n, err := r.server.novelClient.UpdateNovel(c, req)
	if err != nil {
		return nil, handleError("UpdateNovel", err)
	}

	return protoNovelToGraphQL(n), nil
}

func (r *mutationResolver) DeleteNovel(c context.Context, id string) (*DeleteResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	res, err := r.server.novelClient.DeleteNovel(c, id)
	if err != nil {
		return nil, handleError("DeleteNovel", err)
	}

	return &DeleteResponse{DeletedID: res.DeletedId, Success: res.Success, Message: res.Message}, nil
}


func (r *mutationResolver) CreateChapter(c context.Context, in ChapterInput) (*Chapter, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	req := &novelGenproto.CreateChapterRequest{
		NovelId:       in.NovelID,
		ChapterNumber: in.ChapterNumber,
	}
	if in.Title != nil {
		req.Title = *in.Title
	}
	if in.TranslatorGroupID != nil {
		req.TranslatorGroupId = *in.TranslatorGroupID
	}
	if in.SourceURL != nil {
		req.SourceUrl = *in.SourceURL
	}

	ch, err := r.server.novelClient.CreateChapter(c, req)
	if err != nil {
		return nil, handleError("CreateChapter", err)
	}

	return protoChapterToGraphQL(ch), nil
}

func (r *mutationResolver) UpdateChapter(c context.Context, id string, in ChapterInput) (*Chapter, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	req := &novelGenproto.UpdateChapterRequest{
		Id:            id,
		NovelId:       in.NovelID,
		ChapterNumber: in.ChapterNumber,
	}
	if in.Title != nil {
		req.Title = *in.Title
	}
	if in.TranslatorGroupID != nil {
		req.TranslatorGroupId = *in.TranslatorGroupID
	}
	if in.SourceURL != nil {
		req.SourceUrl = *in.SourceURL
	}

	ch, err := r.server.novelClient.UpdateChapter(c, req)
	if err != nil {
		return nil, handleError("UpdateChapter", err)
	}

	return protoChapterToGraphQL(ch), nil
}

func (r *mutationResolver) DeleteChapter(c context.Context, id string) (*DeleteResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	res, err := r.server.novelClient.DeleteChapter(c, id)
	if err != nil {
		return nil, handleError("DeleteChapter", err)
	}

	return &DeleteResponse{DeletedID: res.DeletedId, Success: res.Success, Message: res.Message}, nil
}


func (r *mutationResolver) CreateAuthor(c context.Context, in AuthorInput) (*Author, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	bio := ""
	if in.Bio != nil {
		bio = *in.Bio
	}

	a, err := r.server.novelClient.CreateAuthor(c, in.Name, bio)
	if err != nil {
		return nil, handleError("CreateAuthor", err)
	}

	var createdAt time.Time
	if a.CreatedAt != nil {
		createdAt.UnmarshalBinary(a.CreatedAt)
	}

	return &Author{ID: a.Id, Name: a.Name, Bio: &a.Bio, CreatedAt: createdAt}, nil
}

func (r *mutationResolver) CreateTag(c context.Context, in TagInput) (*Tag, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	slug := ""
	if in.Slug != nil {
		slug = *in.Slug
	}

	t, err := r.server.novelClient.CreateTag(c, in.Name, slug)
	if err != nil {
		return nil, handleError("CreateTag", err)
	}

	return &Tag{ID: int(t.Id), Name: t.Name, Slug: t.Slug}, nil
}

// ── Translation Group ──────────────────

func (r *mutationResolver) CreateTranslationGroup(c context.Context, in TranslationGroupInput) (*TranslationGroup, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	websiteUrl, description := "", ""
	if in.WebsiteURL != nil {
		websiteUrl = *in.WebsiteURL
	}
	if in.Description != nil {
		description = *in.Description
	}

	g, err := r.server.novelClient.CreateTranslationGroup(c, in.Name, websiteUrl, description)
	if err != nil {
		return nil, handleError("CreateTranslationGroup", err)
	}

	var createdAt time.Time
	if g.CreatedAt != nil {
		createdAt.UnmarshalBinary(g.CreatedAt)
	}

	return &TranslationGroup{
		ID: g.Id, Name: g.Name, WebsiteURL: &g.WebsiteUrl,
		Description: &g.Description, CreatedAt: createdAt,
	}, nil
}

// ── Reading List ───────────────────────

func (r *mutationResolver) AddToReadingList(c context.Context, accountID string, in ReadingListInput) (*ReadingListEntry, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	currentChapter := float64(0)
	rating := int32(0)
	notes := ""
	isFavorite := false

	if in.CurrentChapter != nil {
		currentChapter = *in.CurrentChapter
	}
	if in.Rating != nil {
		rating = int32(*in.Rating)
	}
	if in.Notes != nil {
		notes = *in.Notes
	}
	if in.IsFavorite != nil {
		isFavorite = *in.IsFavorite
	}

	e, err := r.server.readinglistClient.AddToReadingList(c, accountID, in.NovelID, in.Status, currentChapter, rating, notes, isFavorite)
	if err != nil {
		return nil, handleError("AddToReadingList", err)
	}

	ratingInt := int(e.Rating)
	return &ReadingListEntry{
		ID: e.ID, NovelID: e.NovelID, Status: e.Status,
		CurrentChapter: e.CurrentChapter, Rating: &ratingInt,
		Notes: &e.Notes, IsFavorite: e.IsFavorite,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}, nil
}

func (r *mutationResolver) UpdateReadingList(c context.Context, id string, in ReadingListInput) (*ReadingListEntry, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	currentChapter := float64(0)
	rating := int32(0)
	notes := ""
	isFavorite := false

	if in.CurrentChapter != nil {
		currentChapter = *in.CurrentChapter
	}
	if in.Rating != nil {
		rating = int32(*in.Rating)
	}
	if in.Notes != nil {
		notes = *in.Notes
	}
	if in.IsFavorite != nil {
		isFavorite = *in.IsFavorite
	}

	e, err := r.server.readinglistClient.UpdateReadingList(c, id, in.Status, currentChapter, rating, notes, isFavorite)
	if err != nil {
		return nil, handleError("UpdateReadingList", err)
	}

	ratingInt := int(e.Rating)
	return &ReadingListEntry{
		ID: e.ID, NovelID: e.NovelID, Status: e.Status,
		CurrentChapter: e.CurrentChapter, Rating: &ratingInt,
		Notes: &e.Notes, IsFavorite: e.IsFavorite,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}, nil
}

func (r *mutationResolver) RemoveFromReadingList(c context.Context, id string) (*DeleteResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	res, err := r.server.readinglistClient.RemoveFromReadingList(c, id)
	if err != nil {
		return nil, handleError("RemoveFromReadingList", err)
	}

	return &DeleteResponse{DeletedID: res.DeletedId, Success: res.Success, Message: res.Message}, nil
}


func (r *mutationResolver) CreateReview(c context.Context, in ReviewInput) (*Review, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	if in.Rating < 1 || in.Rating > 5 {
		return nil, ErrInvalidParameter
	}

	title, content := "", ""
	isSpoiler := false
	if in.Title != nil {
		title = *in.Title
	}
	if in.Content != nil {
		content = *in.Content
	}
	if in.IsSpoiler != nil {
		isSpoiler = *in.IsSpoiler
	}

	rev, err := r.server.reviewClient.PostReview(c, in.NovelID, in.AccountID, title, content, int32(in.Rating), isSpoiler)
	if err != nil {
		return nil, handleError("CreateReview", err)
	}

	return &Review{
		ID: rev.ID, NovelID: rev.NovelID, AccountID: rev.AccountID,
		Rating: rev.Rating, Title: &rev.Title, Content: &rev.Content,
		IsSpoiler: rev.IsSpoiler, CreatedAt: rev.CreatedAt,
	}, nil
}

func (r *mutationResolver) DeleteReview(c context.Context, id string) (*DeleteResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	res, err := r.server.reviewClient.DeleteReview(c, id)
	if err != nil {
		return nil, handleError("DeleteReview", err)
	}

	return &DeleteResponse{DeletedID: res.DeletedId, Success: res.Success, Message: res.Message}, nil
}


func (r *mutationResolver) IncrementViewCount(c context.Context, novelID string) (int, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	count, err := r.server.novelClient.IncrementViewCount(c, novelID)
	if err != nil {
		return 0, handleError("IncrementViewCount", err)
	}

	return int(count), nil
}


func protoNovelToGraphQL(n *novelGenproto.Novel) *Novel {
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

	novel := &Novel{
		ID:            n.Id,
		Title:         n.Title,
		Status:        n.Status,
		NovelType:     n.NovelType,
		TotalChapters: int(n.TotalChapters),
		RatingAvg:     n.RatingAvg,
		RatingCount:   int(n.RatingCount),
		ViewCount:     int(n.ViewCount),
		BookmarkCount: int(n.BookmarkCount),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if n.AlternativeTitle != "" {
		novel.AlternativeTitle = &n.AlternativeTitle
	}
	if n.Description != "" {
		novel.Description = &n.Description
	}
	if n.CoverImageUrl != "" {
		novel.CoverImageURL = &n.CoverImageUrl
	}
	if n.CountryOfOrigin != "" {
		novel.CountryOfOrigin = &n.CountryOfOrigin
	}
	if n.YearPublished != 0 {
		yp := int(n.YearPublished)
		novel.YearPublished = &yp
	}

	if n.Author != nil {
		var authorCreatedAt time.Time
		if n.Author.CreatedAt != nil {
			authorCreatedAt.UnmarshalBinary(n.Author.CreatedAt)
		}
		novel.Author = &Author{ID: n.Author.Id, Name: n.Author.Name, Bio: &n.Author.Bio, CreatedAt: authorCreatedAt}
	}

	for _, g := range n.Genres {
		novel.Genres = append(novel.Genres, &Genre{ID: int(g.Id), Name: g.Name, Slug: g.Slug})
	}
	for _, t := range n.Tags {
		novel.Tags = append(novel.Tags, &Tag{ID: int(t.Id), Name: t.Name, Slug: t.Slug})
	}

	log.Printf("Novel %s (%s) converted", n.Id, n.Title)
	_ = strings.ToLower 

	return novel
}

func protoChapterToGraphQL(ch *novelGenproto.Chapter) *Chapter {
	if ch == nil {
		return nil
	}

	var createdAt, updatedAt time.Time
	if ch.CreatedAt != nil {
		createdAt.UnmarshalBinary(ch.CreatedAt)
	}
	if ch.UpdatedAt != nil {
		updatedAt.UnmarshalBinary(ch.UpdatedAt)
	}

	chapter := &Chapter{
		ID:            ch.Id,
		NovelID:       ch.NovelId,
		ChapterNumber: ch.ChapterNumber,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if ch.Title != "" {
		chapter.Title = &ch.Title
	}
	if ch.TranslatorGroupId != "" {
		chapter.TranslatorGroupID = &ch.TranslatorGroupId
	}
	if ch.SourceUrl != "" {
		chapter.SourceURL = &ch.SourceUrl
	}

	return chapter
}
